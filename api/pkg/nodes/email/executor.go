package email

import (
	"context"
	"fmt"
	"workflow-code-test/api/pkg/mailer"
)

type Inputs struct {
	Name        string  `json:"name"`
	City        string  `json:"city"`
	Email       string  `json:"email"`
	Temperature float64 `json:"temperature"`
}

type Outputs struct {
	EmailSent bool `json:"emailSent"`
}

type Options struct {
	MailClient mailer.Client
}

type Executor struct {
	Opts *Options

	args   map[string]any
	inputs Inputs
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse() error {
	name, ok := e.args["name"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get name where it should string", e.ID())
	}

	email, ok := e.args["email"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get email where it should string", e.ID())
	}

	city, ok := e.args["city"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get city where it should string", e.ID())
	}

	temperature, ok := e.args["temperature"].(float64)
	if !ok {
		return fmt.Errorf("%s: validation failed to get temperature where it should float64", e.ID())
	}

	e.inputs = Inputs{
		Name:        name,
		Email:       email,
		City:        city,
		Temperature: temperature,
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "weather-api"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	err := e.Opts.MailClient.Send(e.inputs.Email, "Weather Alert", fmt.Sprintf("Weather alert for %s! Temperature is %.2f°C!", e.inputs.City, e.inputs.Temperature))
	if err != nil {
		return &Outputs{
			EmailSent: false,
		}, fmt.Errorf("%s: failed to send email: %w", e.ID(), err)
	}

	return &Outputs{
		EmailSent: true,
	}, nil
}
