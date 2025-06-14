package condition_test

import (
	"context"
	"testing"
	"workflow-code-test/api/pkg/nodes/condition"

	"github.com/stretchr/testify/require"
)

func TestConditionExecute(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		threshold   float64
		operator    condition.Operator
		expected    bool
	}{
		{
			name:        "temperature greater than threshold",
			temperature: 28.5,
			threshold:   25,
			operator:    condition.GreaterThanOperator,
			expected:    true,
		},
		{
			name:        "temperature less than threshold",
			temperature: 20.0,
			threshold:   25,
			operator:    condition.LessThanOperator,
			expected:    true,
		},
		{
			name:        "temperature equal to threshold",
			temperature: 25.0,
			threshold:   25,
			operator:    condition.EqualToOperator,
			expected:    true,
		},
		{
			name:        "temperature greater than or equal to threshold (greater)",
			temperature: 26.0,
			threshold:   25,
			operator:    condition.IsAtLeastOperator,
			expected:    true,
		},
		{
			name:        "temperature greater than or equal to threshold (equal)",
			temperature: 25.0,
			threshold:   25,
			operator:    condition.IsAtLeastOperator,
			expected:    true,
		},
		{
			name:        "temperature less than or equal to threshold (less)",
			temperature: 24.0,
			threshold:   25,
			operator:    condition.IsAtMostOperator,
			expected:    true,
		},
		{
			name:        "temperature less than or equal to threshold (equal)",
			temperature: 25.0,
			threshold:   25,
			operator:    condition.IsAtMostOperator,
			expected:    true,
		},
		{
			name:        "temperature not greater than threshold",
			temperature: 24.0,
			threshold:   25,
			operator:    condition.GreaterThanOperator,
			expected:    false,
		},
		{
			name:        "temperature not less than threshold",
			temperature: 26.0,
			threshold:   25,
			operator:    condition.LessThanOperator,
			expected:    false,
		},
		{
			name:        "temperature not equal to threshold",
			temperature: 24.0,
			threshold:   25,
			operator:    condition.EqualToOperator,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &condition.Options{
				Inputs: map[string]interface{}{
					"temperature": tt.temperature,
				},
				Variables: map[string]interface{}{
					"threshold": tt.threshold,
				},
				Operator:   tt.operator,
				Expression: "{{temperature}} {{operator}} {{threshold}}",
			}

			condition := condition.NewExecutor(opts)
			ctx := context.Background()
			outputs, err := condition.Execute(ctx)

			require.NoError(t, err)
			require.NotNil(t, outputs)
			require.Equal(t, tt.expected, outputs.ConditionMet)
		})
	}
}
