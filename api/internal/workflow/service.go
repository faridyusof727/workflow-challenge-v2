package workflow

import (
	"context"
	"time"
)

type Service struct {
	// TODO: Add repository dependency
}

func (s *Service) Execute(ctx context.Context, workflowID string) (*ExecutionResult, error) {
	// TODO: Fetch workflow definition from repository
	// TODO: Execute workflow nodes
	// TODO: Save execution result

	return &ExecutionResult{
		Status:     ExecutionStatusCompleted,
		ExecutedAt: time.Now(),
		Steps:      []any{}, // TODO: set actual steps
	}, nil
}

func NewExecutor() Executor {
	return &Service{}
}
