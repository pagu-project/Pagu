package notification

import (
	"fmt"
	"github.com/pagu-project/Pagu/pkg/notification/zoho"
	"html/template"
)

type EmailSender struct {
	ProviderConfig
}

type IEmailSender interface {
	SendByTemplate(sender string, recipients []string, tmpl *template.Template, data any) error
}

func (e *EmailSender) SendTemplateMail(provider Provider, sender string, recipients []string, tmpl *template.Template, data any) error {
	switch provider {
	case NotificationProviderZapToMail:
		config, ok := e.ProviderConfig.(zoho.ZapToMailerConfig)
		if !ok {
			return fmt.Errorf("unsupported notification provider: %s", provider)
		}
		z := zoho.NewZapToMailer(config)
		return z.SendByTemplate(sender, recipients, tmpl, data)
	default:
		return fmt.Errorf("unsupported notification provider: %s", provider)
	}
}

func LoadMailTemplate(path string) (*template.Template, error) {
	tmpl, err := template.New("").ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
