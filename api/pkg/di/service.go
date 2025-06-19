package di

import (
	"context"
	"os"
	"workflow-code-test/api/pkg/config"
)

type serviceImpl struct {
	container *Container
	config    *config.Config
}

// Config implements Service.
func (s *serviceImpl) Config() *config.Config {
	return s.config
}

// Containers implements Service.
func (s *serviceImpl) Container(ctx context.Context) *Container {
	s.container = &Container{}

	logger := s.logger()
	s.container.Logger = logger

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("unable to parse configs")
		os.Exit(1)
	}
	s.config = cfg

	dbService := s.dbService(ctx, cfg)
	s.container.DbService = dbService

	nodeService := s.nodeService()
	s.container.NodeService = nodeService

	return s.container
}

// Shutdown implements Service.
func (s *serviceImpl) Shutdown(ctx context.Context) error {
	s.container.DbService.Disconnect(ctx)
	return nil
}

func NewService() Service {
	return &serviceImpl{}
}
