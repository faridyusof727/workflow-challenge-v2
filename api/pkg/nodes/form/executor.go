package form

import (
	"context"
	"fmt"
)

type Executor struct {
	args map[string]any
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse() error {
	_, ok := e.args["name"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get name where it should string", e.ID())
	}

	_, ok = e.args["email"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get email where it should string", e.ID())
	}

	_, ok = e.args["city"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get city where it should string", e.ID())
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "form"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	return e.args, nil
}
