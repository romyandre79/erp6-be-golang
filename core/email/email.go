package email

import (
	"fmt"
	"os"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

var registry = map[string]func() (EmailSender, error){}

func Register(name string, factory func() (EmailSender, error)) {
	registry[name] = factory
}

func NewEmailSender() (EmailSender, error) {
	backend := os.Getenv("EMAIL_TYPE") // smtp, sendgrid, mailgun, ses
	if backend != "" {
		if factory, ok := registry[backend]; ok {
			return factory()
		}
		return nil, fmt.Errorf("unknown email backend: %s", backend)
	} else {
		return nil, nil
	}
}
