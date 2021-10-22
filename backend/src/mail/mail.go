package mail

import (
	"app/config"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"embed"
)

//go:embed templates/*
var templates embed.FS

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

type Sender struct {
	Template     string
	To           string
	Subject      string
	Placeholders map[string]string
}

func NewSender(template string, to string, subject string, placeholders map[string]string) *Sender {
	s := &Sender{
		Template:     template,
		To:           to,
		Subject:      subject,
		Placeholders: placeholders,
	}
	return s
}

func (sender *Sender) Send() error {
	subject := makeSubject(sender.Subject)
	body, err := makeBody(sender.Template, sender.Placeholders)
	if err != nil {
		return err
	}

	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			fmt.Sprintf("%s@%s", config.Get().Mail.Smtp.User, config.Get().Mail.Smtp.Host),
			config.Get().Mail.Smtp.Password,
			config.Get().Mail.Smtp.Host,
		),
	}
	from := mail.Address{
		Name:    config.Get().App.Name,
		Address: config.Get().Mail.From,
	}
	msg := []byte(
		fmt.Sprintf("From: %s\r\n", from.String()) +
			fmt.Sprintf("To: %s\r\n", sender.To) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"\r\n" +
			body)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", config.Get().Mail.Smtp.Host, config.Get().Mail.Smtp.Port),
		auth,
		config.Get().Mail.From,
		[]string{sender.To},
		msg)
	return err
}

func makeSubject(subject string) string {
	subject = "[" + config.Get().App.Name + "] " + subject
	return subject
}

func makeBody(template string, placeholders map[string]string) (string, error) {
	b, err := templates.ReadFile(fmt.Sprintf("templates/%s", template))
	if err != nil {
		return "", err
	}

	body := string(b)
	if placeholders != nil {
		for k, v := range placeholders {
			body = strings.Replace(body, k, v, -1)
		}
	}

	// Append Signature.
	sigb, _ := templates.ReadFile(fmt.Sprintf("templates/%s", "signature"))
	sig := string(sigb)
	body = strings.Replace(body, "@SIGNATURE@", sig, -1)

	return body, nil
}
