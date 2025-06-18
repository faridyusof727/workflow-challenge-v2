package di

import (
	"context"
	"log/slog"

	"workflow-code-test/api/pkg/postgres"
)

type Service interface {
	Container(ctx context.Context) *Container
	Shutdown(ctx context.Context) error
}

type Container struct {
	Logger    *slog.Logger
	DbService *postgres.Service
}
