package workflow

import (
	"context"
)

// WorkflowExecutor defines the interface for workflow execution
type Executor interface {
	Execute(ctx context.Context, workflowID string) (*ExecutionResult, error)
}
