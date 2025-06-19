package email

import (
	"context"
	"fmt"
	"workflow-code-test/api/pkg/mailer"
)

type Options struct {
	MailClient mailer.Client
}

type Executor struct {
	Opts *Options

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
	return "email"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	err := e.Opts.MailClient.Send(
		e.args["email"].(string),
		"Weather Alert",
		fmt.Sprintf("Weather alert for %s! Temperature is %s°C!",
			e.args["city"].(string),
			e.args["temperature"].(string),
		),
	)
	if err != nil {
		return map[string]any{
			"emailSent": false,
		}, fmt.Errorf("%s: failed to send email: %w", e.ID(), err)
	}

	// Hardcoded for now to explicitly there should be one output from the mail execution
	if len(e.outputFields) != 1 {
		return nil, fmt.Errorf("%s: output should only contain one variable, outputs: %+v", e.ID(), e.outputFields)
	}

	result := map[string]any{}
	for _, field := range e.outputFields {
		result[field] = true
	}

	return result, nil
}
