package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workflow-code-test/api/pkg/di"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	di *di.Container
}

func NewServer(di *di.Container) *Server {
	return &Server{
		di: di,
	}
}

func (s *Server) Start() {
	container := s.di
	mainRouter := mux.NewRouter()
	mainRouter.Use(RecoverMiddleware(container.Logger))
	mainRouter.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3003"}), // Frontend URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	))

	apiRouter := mainRouter.PathPrefix("/api/v1").Subrouter()

	apiService, err := NewRouter(container)
	if err != nil {
		container.Logger.Error("Failed to create workflow service", "error", err)
		os.Exit(1)
	}

	apiService.LoadRoutes(apiRouter, false)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mainRouter,
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		container.Logger.Info("Starting server on :8080")
		serverErrors <- srv.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking select waiting for either a signal or an error
	select {
	case err := <-serverErrors:
		container.Logger.Error("Server error", "error", err)

	case sig := <-shutdown:
		container.Logger.Info("Shutdown signal received", "signal", sig)

		// Give outstanding requests 5 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			container.Logger.Error("Could not stop server gracefully", "error", err)
			srv.Close()
		}
	}
}
