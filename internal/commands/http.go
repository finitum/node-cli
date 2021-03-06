// Copyright © 2017 The virtual-kubelet authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/finitum/node-cli/internal/stats"
	"github.com/finitum/node-cli/opts"
	"github.com/finitum/node-cli/provider"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
)

// AcceptedCiphers is the list of accepted TLS ciphers, with known weak ciphers elided
// Note this list should be a moving target.
var AcceptedCiphers = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,

	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

func loadTLSConfig(certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "error loading tls certs")
	}

	return &tls.Config{
		Certificates:             []tls.Certificate{cert},
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		CipherSuites:             AcceptedCiphers,
	}, nil
}

func setupHTTPServer(ctx context.Context, p provider.Provider, cfg *apiServerConfig) (_ func(), metricsPort, k8sPort int, retErr error) {
	var closers []io.Closer
	cancel := func() {
		for _, c := range closers {
			c.Close()
		}
	}
	defer func() {
		if retErr != nil {
			cancel()
		}
	}()

	if cfg.CertPath == "" || cfg.KeyPath == "" {
		log.G(ctx).
			WithField("certPath", cfg.CertPath).
			WithField("keyPath", cfg.KeyPath).
			Error("TLS certificates not provided, not setting up pod http server")
	} else {
		tlsCfg, err := loadTLSConfig(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, 0, 0, err
		}
		l, err := tls.Listen("tcp", cfg.Addr, tlsCfg)
		if err != nil {
			return nil, 0, 0, errors.Wrap(err, "error setting up listener for pod http server")
		}

		cfg.Addr = l.Addr().String()

		mux := http.NewServeMux()

		podRoutes := api.PodHandlerConfig{
			RunInContainer:   p.RunInContainer,
			GetContainerLogs: p.GetContainerLogs,
			GetPods:          p.GetPods,
		}
		api.AttachPodRoutes(podRoutes, mux, true)

		s := &http.Server{
			Handler:   mux,
			TLSConfig: tlsCfg,
		}
		go serveHTTP(ctx, s, l, "pods")
		closers = append(closers, s)
		k8sPort = l.Addr().(*net.TCPAddr).Port
	}

	if mp, ok := p.(provider.PodMetricsProvider); ok {
		if cfg.MetricsAddr == "" {
			log.G(ctx).Info("Pod metrics server not setup due to empty metrics address")
		} else {
			l, err := net.Listen("tcp", cfg.MetricsAddr)
			if err != nil {
				return nil, 0, 0, errors.Wrap(err, "could not setup listener for pod metrics http server")
			}
			cfg.MetricsAddr = l.Addr().String()
			mux := http.NewServeMux()

			reg := prometheus.NewRegistry()
			var gatherer prometheus.Gatherer = reg
			var registerer prometheus.Registerer = reg

			registerer.MustRegister(stats.NewCollector(&mp))
			mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
				reg, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}),
			))

			s := &http.Server{
				Handler: mux,
			}
			go serveHTTP(ctx, s, l, "pod metrics")
			closers = append(closers, s)
			metricsPort = l.Addr().(*net.TCPAddr).Port
		}
	}

	return cancel, metricsPort, k8sPort, nil
}

func serveHTTP(ctx context.Context, s *http.Server, l net.Listener, name string) {
	if err := s.Serve(l); err != nil {
		select {
		case <-ctx.Done():
			//
		default:
			log.G(ctx).WithError(err).Errorf("Error setting up %s http server", name)
		}
	}
	_ = l.Close()
}

type apiServerConfig struct {
	CertPath    string
	KeyPath     string
	Addr        string
	MetricsAddr string
}

func getAPIConfig(c *opts.Opts) (*apiServerConfig, error) {
	config := apiServerConfig{
		CertPath: os.Getenv("APISERVER_CERT_LOCATION"),
		KeyPath:  os.Getenv("APISERVER_KEY_LOCATION"),
	}

	config.Addr = fmt.Sprintf(":%d", c.ListenPort)
	config.MetricsAddr = c.MetricsAddr

	return &config, nil
}
