package di

import (
	"context"
	"os"
	"workflow-code-test/api/pkg/postgres"
)

func (s *serviceImpl) dbService(ctx context.Context) *postgres.Service {
	pool, err := postgres.NewService(ctx, &postgres.Options{
		ConnectionURI: os.Getenv("DATABASE_URL"),
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
