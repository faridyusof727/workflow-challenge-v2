package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"workflow-code-test/api/pkg/render"

	"github.com/gorilla/mux"
)

// JsonMiddleware sets the Content-Type header to application/json
func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RecoverMiddleware recovers from panics and logs them using the provided logger.
func RecoverMiddleware(log *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log = log.With("stack", string(debug.Stack()))
					render.Error(w, r, http.StatusInternalServerError, render.ErrInternalServerError, log)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
