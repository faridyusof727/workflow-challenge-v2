package condition

import (
	"context"
	"fmt"
	"maps"
	"workflow-code-test/api/pkg/nodes/types"

	"github.com/expr-lang/expr"
)

type Options struct {
	Inputs    map[string]interface{}
	Variables map[string]interface{}
	Operator  Operator

	// Expression represents the condition expression to be evaluated
	// Variables to be used in the condition expression
	//
	// 	Example: {{temperature}} {{operator}} {{threshold}}
	Expression string
}

type Outputs struct {
	ConditionMet bool `json:"conditionMet"`
}

type Executor struct {
	opts *Options
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) Validate() error {
	panic("unimplemented")
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "condition"
}

func NewExecutor(opts *Options) types.NodeExecutor {
	return &Executor{opts: opts}
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	maps.Insert(e.opts.Variables, maps.All(e.opts.Inputs))

	exprString := ExprReplacePlaceholderByMap(e.opts.Expression, e.opts.Variables, e.opts.Operator)

	program, err := expr.Compile(exprString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to compile expression: %w", e.ID(), err)
	}

	output, err := expr.Run(program, e.opts.Variables)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to run expression: %w", e.ID(), err)
	}

	o, ok := output.(bool)
	if !ok {
		return nil, fmt.Errorf("%s: failed to get output: %w", e.ID(), err)
	}

	return &Outputs{
		ConditionMet: o,
	}, nil
}
