package di

import (
	"context"
	"os"
	"workflow-code-test/api/pkg/config"
	"workflow-code-test/api/pkg/postgres"
)

func (s *serviceImpl) dbService(ctx context.Context, cfg *config.Config) *postgres.Service {
	if cfg == nil {
		s.container.Logger.Error("Failed to get config")
		os.Exit(1)
	}

	pool, err := postgres.NewService(ctx, &postgres.Options{
		ConnectionURI: cfg.Database.URL,
	})
	if err != nil {
		s.container.Logger.Error("Failed to create postgres service", "error", err)
		os.Exit(1)
	}

	if err := pool.PoolConn(ctx); err != nil {
		s.container.Logger.Error("Failed to acquire database connection", "error", err)
		os.Exit(1)
	}

	return pool
}
