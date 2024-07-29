package notification

import (
	"fmt"
	"html/template"
)

type NotificationType string

const (
	NotificationTypeMail NotificationType = "mail"
	//other types like sms, firebase,...
)

type Provider string
type ProviderConfig any

const (
	NotificationProviderZapToMail = "zoho"
)

type ISender interface {
	SendTemplateMail(provider Provider, sender string, recipients []string, tmpl *template.Template, data any) error
}

func New(notificationType NotificationType, configs ProviderConfig) (ISender, error) {
	switch notificationType {
	case NotificationTypeMail:
		return &EmailSender{configs}, nil
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", notificationType)
	}
}
