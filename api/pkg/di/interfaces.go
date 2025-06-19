package di

import (
	"context"
	"log/slog"

	"workflow-code-test/api/pkg/config"
	"workflow-code-test/api/pkg/postgres"
)

// Service defines the interface for the dependency injection service.
// It provides access to the container holding all dependencies and application configuration.
type Service interface {
	// Container returns the dependency container for the given context.
	Container(ctx context.Context) *Container
	// Config returns the application configuration.
	Config() *config.Config
	// Shutdown performs a graceful shutdown of the service.
	Shutdown(ctx context.Context) error
}

// Container holds all the dependencies required by the application.
// It provides centralized access to services like logging and database connections.
type Container struct {
	// Logger is the application's structured logger.
	Logger    *slog.Logger
	// DbService provides access to the Postgres database.
	DbService *postgres.Service
}