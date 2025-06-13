package workflow

import (
	"context"
)

// Executor defines the interface for executing a workflow with a given workflow ID.
// It takes a context and workflow ID, and returns an execution result or an error.
type Executor interface {
	Execute(ctx context.Context, workflowID string) (*ExecutionResult, error)
}
