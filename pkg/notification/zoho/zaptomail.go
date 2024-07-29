package zoho

import (
	"bytes"
	"fmt"
	"github.com/go-mail/mail/v2"
	"html/template"
	"time"
)

type ZapToMailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type ZapToMailer struct {
	dialer *mail.Dialer
}

func NewZapToMailer(config ZapToMailerConfig) *ZapToMailer {
	dialer := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	dialer.Timeout = 5 * time.Second

	return &ZapToMailer{
		dialer: dialer,
	}
}

func (z *ZapToMailer) SendByTemplate(sender string, recipients []string, tmpl *template.Template, data any) error {
	subject := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return fmt.Errorf("error executing template with subject: %v", err)
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return fmt.Errorf("error executing plainbody: %v", err)
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return fmt.Errorf("error executing HTML body: %v", err)
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", recipients...)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	//return z.dialer.DialAndSend(msg)
	err = z.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
