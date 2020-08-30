package mailer

import (
	"github.com/carlosstrand/mailer/types"
)

type Mailer struct {
	config *types.MailerConfig
}

func NewMailer(config *types.MailerConfig) *Mailer {
	if len(config.Providers) == 0 {
		panic("You should add at least one mail provider")
	}
	return &Mailer{config: config}
}
