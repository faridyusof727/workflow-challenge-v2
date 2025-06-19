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

	args map[string]any
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse() error {
	_, ok := e.args["name"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get name where it should string", e.ID())
	}

	_, ok = e.args["email"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get email where it should string", e.ID())
	}

	_, ok = e.args["city"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get city where it should string", e.ID())
	}

	_, ok = e.args["temperature"].(float64)
	if !ok {
		return fmt.Errorf("%s: validation failed to get temperature where it should float64", e.ID())
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
		fmt.Sprintf("Weather alert for %s! Temperature is %.2f°C!",
			e.args["city"].(string),
			e.args["temperature"].(float64),
		),
	)
	if err != nil {
		return map[string]any{
			"emailSent": false,
		}, fmt.Errorf("%s: failed to send email: %w", e.ID(), err)
	}

	return map[string]any{
		"emailSent": true,
	}, nil
}
