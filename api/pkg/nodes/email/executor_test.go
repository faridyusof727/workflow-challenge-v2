package email_test

import (
	"context"
	"testing"
	"workflow-code-test/api/pkg/mailer"
	"workflow-code-test/api/pkg/nodes/email"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	executor := &email.Executor{
		Opts: &email.Options{
			MailClient: mailer.NewNoopClient(),
		},
	}
	executor.SetArgs(map[string]any{
		"name":        "John Doe",
		"email":       "johndoe@example.com",
		"city":        "New York",
		"temperature": 25.5,
	})
	err := executor.ValidateAndParse()
	require.NoError(t, err)

	outputs, err := executor.Execute(context.Background())
	require.NoError(t, err)

	expectedOutputs := &email.Outputs{
		EmailSent: true,
	}
	require.Equal(t, expectedOutputs, outputs)
}
