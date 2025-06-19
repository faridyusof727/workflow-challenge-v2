package email

import (
	"context"
	"fmt"
	"workflow-code-test/api/pkg/mailer"
)

const (
	TemplateKey string = "emailTemplate"
)

type Options struct {
	MailClient mailer.Client
}

type Executor struct {
	Opts *Options

	args         map[string]any
	tmpl         map[string]any
	outputFields []string
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

func (e *Executor) SetOutputFields(fields []string) {
	e.outputFields = fields
}

// ValidateAndParse implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse(argsCheck []string) error {
	for _, key := range argsCheck {
		_, ok := e.args[key].(string)
		if !ok {
			return fmt.Errorf("%s: validation key failed, key: %v", e.ID(), key)
		}
	}

	tmpl, ok := e.args[TemplateKey].(map[string]any)
	if !ok {
		return fmt.Errorf("%s: validation failed to get emailTemplate where it should a map", e.ID())
	}

	_, ok = tmpl["body"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get emailTemplate.body where it should a string", e.ID())
	}

	_, ok = tmpl["subject"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get emailTemplate.subject where it should a string", e.ID())
	}

	e.tmpl = tmpl
	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "email"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	err := e.Opts.MailClient.Send(
		e.args["email"].(string),
		e.tmpl["subject"].(string),
		TmplReplacePlaceholderByMap(
			e.tmpl["body"].(string),
			e.args,
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
