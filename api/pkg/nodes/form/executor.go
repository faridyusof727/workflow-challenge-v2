package form

import (
	"context"
	"fmt"
)

// Input has strict fields for now.
// This can be made more flexible in the future.
type Inputs struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	City  string `json:"city"`
}

// Output has strict fields for now.
// This can be made more flexible in the future.
type Outputs struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	City  string `json:"city"`
}

type Executor struct {
	args   map[string]any
	inputs Inputs
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse() error {
	name, ok := e.args["name"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get name where it should string", e.ID())
	}

	email, ok := e.args["email"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get email where it should string", e.ID())
	}

	city, ok := e.args["city"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get city where it should string", e.ID())
	}

	e.inputs = Inputs{
		Name:  name,
		Email: email,
		City:  city,
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "form"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	return &Outputs{
		Name:  e.inputs.Name,
		Email: e.inputs.Email,
		City:  e.inputs.City,
	}, nil
}
