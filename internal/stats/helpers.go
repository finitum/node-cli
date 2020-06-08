package stats

import (
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"io"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func handleError(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := f(w, req)
		if err == nil {
			return
		}

		code := httpStatusCode(err)
		w.WriteHeader(code)
		_, _ = io.WriteString(w, err.Error()) //nolint:errcheck
	}
}

func httpStatusCode(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errdefs.IsNotFound(err):
		return http.StatusNotFound
	case errdefs.IsInvalidInput(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
