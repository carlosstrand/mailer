package types

import "github.com/carlosstrand/mailer/provider"

type MailerConfig struct {
	DefaultFrom string
	PublicPath  string
	Providers   []provider.MailProvider
}

type SendMailFromTemplateRequest struct {
	From     provider.EmailAccount `json:"from"`
	To       provider.EmailAccount `json:"to"`
	Template string                `json:"template"`
	Vars     map[string]string     `json:"vars"`
}
