package types

import (
	"github.com/carlosstrand/mailer/provider"
)

type Template struct {
	Subject string
}

type TemplateMap map[string]Template

type MailerConfig struct {
	DefaultFrom string
	PublicPath  string
	Providers   []provider.MailProvider
	Templates map[string]Template
}

type SendMailFromTemplateRequest struct {
	From     provider.EmailAccount `json:"from"`
	To       provider.EmailAccount `json:"to"`
	Template string                `json:"template"`
	Vars     map[string]string     `json:"vars"`
}
