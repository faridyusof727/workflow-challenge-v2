package condition

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
)

const (
	ExpressionKey string = "conditionExpression"
	OperatorKey   string = "operator"
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

func (e *Executor) ValidateAndParse(argsCheck []string) error {
	_, ok := e.args[ExpressionKey].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get expression where it should string", e.ID())
	}

	operator, ok := e.args[OperatorKey].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get operator where it should &Operator", e.ID())
	}

	if err := Operator(operator).Validate(); err != nil {
		return fmt.Errorf("%s: validation failed to validate operator: %w", e.ID(), err)
	}

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
	return "condition"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	rawExprString := e.args[ExpressionKey].(string)

	exprString := ExprReplacePlaceholderByMap(rawExprString, e.args)

	program, err := expr.Compile(exprString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to compile expression: %w", e.ID(), err)
	}

	output, err := expr.Run(program, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to run expression: %w", e.ID(), err)
	}

	o, ok := output.(bool)
	if !ok {
		return nil, fmt.Errorf("%s: failed to get output: %w", e.ID(), err)
	}

	// Hardcoded for now to explicitly there should be one output from the expression
	if len(e.outputFields) != 1 {
		return nil, fmt.Errorf("%s: output should only contain one variable, outputs: %+v", e.ID(), e.outputFields)
	}

	result := map[string]any{}
	for _, field := range e.outputFields {
		result[field] = o
	}

	return result, nil
}
