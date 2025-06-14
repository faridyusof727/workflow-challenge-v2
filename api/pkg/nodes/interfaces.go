package nodes

import "context"

type NodeExecutor interface {
	// Execute runs the executor's primary logic and returns a map of outputs or an error.
	// The context allows for potential cancellation or timeout of the execution.
	Execute(ctx context.Context) (map[string]interface{}, error)

	// ID returns a unique identifier for the node.
	// This method is typically used to distinguish and reference individual nodes within a system.
	ID() string
}
