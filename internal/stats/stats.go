package stats

import (
	"context"
	"encoding/json"
	"github.com/finitum/node-cli/stats"
	"github.com/pkg/errors"
	"github.com/virtual-kubelet/virtual-kubelet/node/api"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"net/http"
)

// PodStatsSummaryHandlerFunc defines the handler for getting pod stats summaries
type PodStatsSummaryHandlerFunc func(context.Context) (*stats.Summary, error)

// HandlePodStatsSummary makes an HTTP handler for implementing the kubelet summary stats endpoint
func HandlePodStatsSummary(h PodStatsSummaryHandlerFunc) http.HandlerFunc {
	if h == nil {
		return api.NotImplemented
	}
	return handleError(func(w http.ResponseWriter, req *http.Request) error {
		stats, err := h(req.Context())
		if err != nil {
			if isCancelled(err) {
				return err
			}
			return errors.Wrap(err, "error getting status from provider")
		}

		b, err := json.Marshal(stats)
		if err != nil {
			return errors.Wrap(err, "error marshalling stats")
		}

		if _, err := w.Write(b); err != nil {
			return errors.Wrap(err, "could not write to client")
		}
		return nil
	})
}

func isCancelled(err error) bool {
	if err == context.Canceled {
		return true
	}

	if e, ok := err.(causal); ok {
		return isCancelled(e.Cause())
	}
	return false
}

type causal interface {
	Cause() error
	error
}

// InstrumentHandler wraps an http.Handler and injects instrumentation into the request context.
func InstrumentHandler(h http.Handler) http.Handler {
	instrumented := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
	return &ochttp.Handler{
		Handler:     instrumented,
		Propagation: &b3.HTTPFormat{},
	}
}
