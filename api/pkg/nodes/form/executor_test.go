package form_test

import (
	"context"
	"testing"
	"workflow-code-test/api/pkg/nodes/form"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	executor := &form.Executor{}
	executor.SetArgs(map[string]any{
		"name":  "John Doe",
		"email": "johndoe@example.com",
		"city":  "New York",
	})
	executor.SetOutputFields([]string{"name", "email", "city"})
	err := executor.ValidateAndParse([]string{"name", "email", "city"})
	require.NoError(t, err)

	outputs, err := executor.Execute(context.Background())
	require.NoError(t, err)

	expectedOutputs := map[string]any{
		"name":  "John Doe",
		"email": "johndoe@example.com",
		"city":  "New York",
	}
	require.Equal(t, expectedOutputs, outputs)
}
