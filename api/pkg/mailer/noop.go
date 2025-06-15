package mailer

import "log/slog"

type noopClient struct{}

// Send implements Client.
func (n *noopClient) Send(to string, subject string, body string) error {
	slog.Info("noop: email sent!",
		slog.String("to", to),
		slog.String("subject", subject),
		slog.String("body", body),
	)

	return nil
}

func NewNoopClient() Client {
	return &noopClient{}
}
