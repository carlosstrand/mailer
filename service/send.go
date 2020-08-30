package service

import (
	"github.com/carlosstrand/mailer/provider"
	"github.com/carlosstrand/mailer/types"
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"sync"
)

func (m *MailerService) SendFromReq(sReq types.SendMailFromTemplateRequest) error {
	data := sync.Map{}
	for key, value := range sReq.Vars {
		data.Store(key, value)
	}
	html, err := m.RenderToString(sReq.Template, &data, true, true)
	if err != nil {
		return err
	}
	message := provider.Message{
		From:      sReq.From,
		To:        sReq.To,
		Subject:   "Mensagem de teste API SendGrid",
		PlainText: "Mensagem testando",
		HTML:      html,
	}
	return m.Send(message)
}

func (m *MailerService) Send(message provider.Message) error {
	for _, p := range m.config.Providers {
		color.Cyan("Sending mail using %s provider...\n", p.Name())
		err := p.Send(message)
		if err == nil {
			return nil
		}
		color.Red("Failed to send message with %s provider. Attempting next provider...", p.Name())
	}
	return nil
}
