package form

import (
	"context"
	"fmt"
)

type Executor struct {
	args         map[string]any
	outputFields []string
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

func (e *Executor) SetOutputFields(fields []string) {
	e.outputFields = fields
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse(argsCheck []string) error {
	for _, key := range argsCheck {
		_, ok := e.args[key].(string)
		if !ok {
			return fmt.Errorf("%s: validation key failed, key: %v", e.ID(), key)
		}
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "form"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	result := map[string]any{}
	for _, field := range e.outputFields {
		if val, ok := e.args[field]; ok {
			result[field] = val
		}

	}

	return result, nil
}
