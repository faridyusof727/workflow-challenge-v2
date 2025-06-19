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
		"temperature": "25.5",
		"emailTemplate": map[string]any{
			"body":    "Weather alert for {{city}}! Temperature is {{temperature}}°C!",
			"subject": "Weather Alert",
		},
	})
	executor.SetOutputFields([]string{"emailSent"})
	err := executor.ValidateAndParse([]string{"name", "email", "city", "temperature"})
	require.NoError(t, err)

	outputs, err := executor.Execute(context.Background())
	require.NoError(t, err)

	result := outputs.(map[string]any)
	require.Equal(t, true, result["emailSent"])
}
