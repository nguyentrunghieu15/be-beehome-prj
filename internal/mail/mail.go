package mail

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type ContentType string

const (
	PLAIN_TEXT ContentType = "text/plain"
	HTML       ContentType = "text/html"
)

type Letter struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Type    ContentType
	Body    string
}

type AuthStmp struct {
	Email    string
	Password string
}

func sendMail(auth AuthStmp, letter Letter) error {
	m := gomail.NewMessage()
	m.SetHeader("From", letter.From)
	m.SetHeader("To", letter.To...)
	m.SetHeader("Cc", letter.Cc...)
	m.SetHeader("Subject", letter.Subject)
	m.SetBody(string(letter.Type), letter.Body)

	host := os.Getenv("MAIL_HOST")
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))

	if err != nil {
		return err
	}

	d := gomail.NewDialer(host, port, auth.Email, auth.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

type IMailBox interface {
	SendMail(AuthStmp, Letter) error
}
type MailBox struct{}

func (*MailBox) SendMail(auth AuthStmp, letter Letter) error {
	return sendMail(auth, letter)
}
