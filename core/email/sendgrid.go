package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridSender struct {
	client *sendgrid.Client
	from   string
}

func init() {
	Register("sendgrid", func() (EmailSender, error) {
		return NewSendGridSender()
	})
}

func NewSendGridSender() (*SendGridSender, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SENDGRID_API_KEY missing")
	}
	from := os.Getenv("EMAIL_FROM")
	return &SendGridSender{
		client: sendgrid.NewSendClient(apiKey),
		from:   from,
	}, nil
}

func (s *SendGridSender) Send(to, subject, body string) error {
	from := mail.NewEmail("App", s.from)
	toAddr := mail.NewEmail("User", to)
	message := mail.NewSingleEmail(from, subject, toAddr, body, body)
	_, err := s.client.Send(message)
	return err
}
