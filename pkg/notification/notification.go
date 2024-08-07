package notification

import (
	"fmt"
	"html/template"
)

type (
	Provider         string
	ProviderConfig   any
	NotificationType int
)

const (
	NotificationTypeMail NotificationType = 0
	// other types like sms, firebase,...
)

func (n NotificationType) String() string {
	switch n {
	case NotificationTypeMail:
		return "mail"
	default:
		return ""
	}
}

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
		return nil, fmt.Errorf("unsupported notification type: %s", notificationType.String())
	}
}
