package email

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESSender struct {
	client *ses.SES
	from   string
}

func init() {
	Register("ses", func() (EmailSender, error) {
		return NewSESSender()
	})
}

func NewSESSender() (*SESSender, error) {
	region := os.Getenv("AWS_REGION")
	from := os.Getenv("EMAIL_FROM")

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, err
	}

	return &SESSender{
		client: ses.New(sess),
		from:   from,
	}, nil
}

func (s *SESSender) Send(to, subject, body string) error {
	input := &ses.SendEmailInput{
		Source: aws.String(s.from),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Subject: &ses.Content{Data: aws.String(subject)},
			Body: &ses.Body{
				Text: &ses.Content{Data: aws.String(body)},
			},
		},
	}
	_, err := s.client.SendEmail(input)
	return err
}
