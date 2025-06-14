package nodes

import "context"

type NodeExecutor[T any] interface {
	// Execute runs the executor's primary logic and returns a map of outputs or an error.
	// The context allows for potential cancellation or timeout of the execution.
	Execute(ctx context.Context) (*T, error)

	// ID returns a unique identifier for the node.
	// This method is typically used to distinguish and reference individual nodes within a system.
	ID() string

	// Validate checks the validity of the node or its configuration.
	// It returns an error if the node is not in a valid state, otherwise returns nil.
	Validate() error
}
