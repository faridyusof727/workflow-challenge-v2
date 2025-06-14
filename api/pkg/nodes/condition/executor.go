package condition

import (
	"context"
	"fmt"
	"maps"
	"workflow-code-test/api/pkg/nodes"

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

type Executor struct {
	opts *Options
}

// ID implements NodeExecutor.
func (c *Executor) ID() string {
	return "condition"
}

func NewExecutor(opts *Options) nodes.NodeExecutor {
	return &Executor{opts: opts}
}

func (c *Executor) Execute(ctx context.Context) (map[string]interface{}, error) {
	maps.Insert(c.opts.Variables, maps.All(c.opts.Inputs))

	exprString := ExprReplacePlaceholderByMap(c.opts.Expression, c.opts.Variables, c.opts.Operator)

	program, err := expr.Compile(exprString)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	output, err := expr.Run(program, c.opts.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to run expression: %w", err)
	}

	return map[string]interface{}{
		"conditionMet": output,
	}, nil
}
