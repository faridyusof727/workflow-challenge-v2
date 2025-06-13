package node

import (
	"context"
)

type Service struct {
	// TODO: Add repository dependency
}

// Execute implements Executor.
func (s *Service) Execute(ctx context.Context, nodeID string) (*ExecutionResult, error) {
	panic("unimplemented")
}

func NewExecutor() Executor {
	return &Service{}
}
