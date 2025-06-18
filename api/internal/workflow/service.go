package workflow

import (
	"context"
	"time"
)

type ServiceImpl struct {
	repo Repository
}

// Workflow implements Service.
func (s *ServiceImpl) Workflow(ctx context.Context, workflowID string) (*Workflow, error) {
	workflow, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *ServiceImpl) Execute(ctx context.Context, workflowID string) (*ExecutionResult, error) {
	s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)

	return &ExecutionResult{
		Status:     ExecutionStatusCompleted,
		ExecutedAt: time.Now(),
		Steps:      []Step{}, // TODO: set actual steps
	}, nil
}

func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}
