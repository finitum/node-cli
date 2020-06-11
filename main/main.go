package main

import (
	"context"
	"github.com/finitum/node-cli/internal/stats"
	"github.com/finitum/node-cli/provider"
	types "github.com/finitum/node-cli/stats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
)

type Test struct{}

func (t *Test) GetStatsSummary() (*types.Summary, error) {
	return &types.Summary{
		Node: types.NodeStats{
			Name:                 "apatelet-8787897123",
			UsageNanoCores:       uint64(42),
			AvailableBytesMemory: uint64(42),
			UsageBytesMemory:     uint64(42),

			AvailableBytesEphemeral: uint64(42),
			CapacityBytesEphemeral:  uint64(42),
			UsedBytesEphemeral:      uint64(42),

			AvailableBytesStorage: uint64(42),
			CapacityBytesStorage:  uint64(42),
			UsedBytesStorage:      uint64(42),
		},
		Pods: []types.PodStats{
			{
				PodRef: types.PodReference{
					Name:      "test_name",
					Namespace: "test_namespace",
					UID:       "test_uid",
				},

				UsageNanoCores: uint64(42),

				AvailableBytesMemory: uint64(42),
				UsageBytesMemory:     uint64(42),

				AvailableBytesEphemeral: uint64(42),
				CapacityBytesEphemeral:  uint64(42),
				UsedBytesEphemeral:      uint64(42),

				AvailableBytesStorage: uint64(42),
				CapacityBytesStorage:  uint64(42),
				UsedBytesStorage:      uint64(42),
			},
			{
				PodRef: types.PodReference{
					Name:      "test_name2",
					Namespace: "test_namespace",
					UID:       "test_uid2",
				},

				UsageNanoCores: uint64(41),

				AvailableBytesMemory: uint64(41),
				UsageBytesMemory:     uint64(41),

				AvailableBytesEphemeral: uint64(41),
				CapacityBytesEphemeral:  uint64(41),
				UsedBytesEphemeral:      uint64(41),

				AvailableBytesStorage: uint64(41),
				CapacityBytesStorage:  uint64(41),
				UsedBytesStorage:      uint64(41),
			},
		},
	}, nil
}

func T() provider.PodMetricsProvider {
	return &Test{}
}

func main() {
	ctx := context.Background()
	l, err := net.Listen("tcp", ":8080")
	log.Printf("%v\n", err)
	mux := http.NewServeMux()
	mp := T()

	prometheus.MustRegister(stats.NewCollector(&mp))

	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Handler: mux,
	}

	serveHTTP(ctx, s, l, "pod metrics")
}

func serveHTTP(ctx context.Context, s *http.Server, l net.Listener, name string) {
	if err := s.Serve(l); err != nil {
		select {
		case <-ctx.Done():
			//
		default:
			//
		}
	}
	_ = l.Close()
}
