package mailer

type Client interface {
	// Send sends an email to the specified recipient with the given subject and body.
	// Returns an error if the email cannot be sent.
	Send(to, subject, body string) error
}
