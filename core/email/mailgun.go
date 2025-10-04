package email

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunSender struct {
	client *mailgun.MailgunImpl
	domain string
	from   string
}

func init() {
	Register("mailgun", func() (EmailSender, error) {
		return NewMailgunSender()
	})
}

func NewMailgunSender() (*MailgunSender, error) {
	domain := os.Getenv("MAILGUN_DOMAIN")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	from := os.Getenv("EMAIL_FROM")

	if domain == "" || apiKey == "" {
		return nil, fmt.Errorf("Mailgun config missing")
	}

	mg := mailgun.NewMailgun(domain, apiKey)
	return &MailgunSender{client: mg, domain: domain, from: from}, nil
}

func (m *MailgunSender) Send(to, subject, body string) error {
	msg := m.client.NewMessage(m.from, subject, body, to)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, _, err := m.client.Send(ctx, msg)
	return err
}
