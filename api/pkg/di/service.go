package di

import (
	"context"
)

type serviceImpl struct {
	container *Container
}

// Containers implements Service.
func (s *serviceImpl) Container(ctx context.Context) *Container {
	s.container = &Container{}

	logger := s.logger()
	s.container.Logger = logger

	dbService := s.dbService(ctx)
	s.container.DbService = dbService

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
