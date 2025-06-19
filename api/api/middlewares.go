package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"workflow-code-test/api/pkg/config"
	"workflow-code-test/api/pkg/render"

	"github.com/gorilla/handlers"
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

// CorsMiddleware creates a middleware function for handling Cross-Origin Resource Sharing (CORS) policies.
// It configures CORS based on the provided configuration and returns a middleware function
// that can be added to a Gorilla Mux router.
//
// Parameters:
//   - cfg: A pointer to the application configuration containing CORS settings.
//
// Returns:
//   - A mux.MiddlewareFunc that can be used with a Gorilla Mux router to apply CORS policies.
func CorsMiddleware(cfg *config.Config) mux.MiddlewareFunc {
	return handlers.CORS(
		handlers.AllowedOrigins(cfg.CORS.AllowedOrigins), // Frontend URL
		handlers.AllowedMethods(cfg.CORS.AllowedMethods),
		handlers.AllowedHeaders(cfg.CORS.AllowedHeaders),
		handlers.AllowCredentials(),
	)
}
