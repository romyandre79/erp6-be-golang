package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type SmtpSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func init() {
	Register("smtp", func() (EmailSender, error) {
		return NewSmtpSender(), nil
	})
}

func NewSmtpSender() *SmtpSender {
	return &SmtpSender{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (s *SmtpSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.from, to, subject, body))
	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{to}, msg)
}
