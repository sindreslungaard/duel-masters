package internal

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sirupsen/logrus"
)

var svc *ses.SES

func SendMail(to string, subject string, body string) error {

	if svc == nil {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String("eu-west-1"),
			Credentials: credentials.NewStaticCredentials(os.Getenv("aws_key"), os.Getenv("aws_secret"), ""),
		})

		if err != nil {
			return err
		}

		svc = ses.New(sess)
	}

	_, err := svc.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String("\"Shobu.io\" <noreply@shobu.io>"),
	})

	if err != nil {
		logrus.Error("Error sending email", err)
		return err
	}

	logrus.Info("Sent an email to %s with subject %s", to, subject)

	return nil

}
