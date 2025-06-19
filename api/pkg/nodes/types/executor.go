package types

import "context"

type NodeExecutor interface {
	// Execute runs the executor's primary logic and returns a map of outputs or an error.
	// The context allows for potential cancellation or timeout of the execution.
	Execute(ctx context.Context) (any, error)

	// ID returns a unique identifier for the node.
	// This method is typically used to distinguish and reference individual nodes within a system.
	ID() string

	// SetArgs configures the node with input arguments.
	// These arguments can be used to customize the node's behaviour or provide input data
	// for its execution. The method allows setting multiple arguments via a map.
	SetArgs(args map[string]any)

	// SetOutputFields specifies which fields should be included in the node's output.
	// args contains the field names that the node should return as part of its execution result.
	// These fields can reference keys from the input arguments set via SetArgs.
	SetOutputFields(fields []string)

	// ValidateAndParse checks the integrity and parses the node's configuration or input arguments.
	// argsCheck specifies which argument names should be validated during the process.
	// Returns an error if validation fails or parsing encounters issues.
	ValidateAndParse(argsCheck []string) error
}
