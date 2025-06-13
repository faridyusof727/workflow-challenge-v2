package node

import (
	"context"
)

// Executor defines an interface for executing a node with a given context and node ID.
// It returns an ExecutionResult and an optional error if the execution fails.
type Executor interface {
	Execute(ctx context.Context, nodeID string) (*ExecutionResult, error)
}
