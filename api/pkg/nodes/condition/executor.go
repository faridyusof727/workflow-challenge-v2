package condition

import (
	"context"
	"fmt"

	"github.com/expr-lang/expr"
)

type Inputs struct {
	Expression  string   `json:"expression"`
	Threshold   float64  `json:"threshold"`
	Operator    Operator `json:"operator"`
	Temperature float64  `json:"temperature"`
}

type Executor struct {
	args   map[string]any
	inputs Inputs
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

func (e *Executor) ValidateAndParse() error {
	expression, ok := e.args["expression"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get expression where it should string", e.ID())
	}

	threshold, ok := e.args["threshold"].(float64)
	if !ok {
		return fmt.Errorf("%s: validation failed to get threshold where it should float64", e.ID())
	}

	operator, ok := e.args["operator"].(Operator)
	if !ok {
		return fmt.Errorf("%s: validation failed to get operator where it should &Operator", e.ID())
	}
	if err := operator.Validate(); err != nil {
		return fmt.Errorf("%s: validation failed to validate operator: %w", e.ID(), err)
	}

	temperature, ok := e.args["temperature"].(float64)
	if !ok {
		return fmt.Errorf("%s: validation failed to get temperature where it should float64", e.ID())
	}

	e.inputs = Inputs{
		Expression:  expression,
		Threshold:   threshold,
		Operator:    operator,
		Temperature: temperature,
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "condition"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	exprString := ExprReplacePlaceholderByMap(e.inputs)

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

	return map[string]any{
		"conditionMet": o,
	}, nil
}
