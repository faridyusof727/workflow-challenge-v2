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
	err := executor.ValidateAndParse()
	require.NoError(t, err)

	outputs, err := executor.Execute(context.Background())
	require.NoError(t, err)

	expectedOutputs := &form.Outputs{
		Name:  "John Doe",
		Email: "johndoe@example.com",
		City:  "New York",
	}
	require.Equal(t, expectedOutputs, outputs)
}
