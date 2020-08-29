package service

import (
	"github.com/carlosstrand/mailer/provider"
	"github.com/carlosstrand/mailer/types"
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"sync"
)

func (s *MailerService) SendFromReq(sReq types.SendMailFromTemplateRequest) error {
	data := sync.Map{}
	for key, value := range sReq.Vars {
		data.Store(key, value)
	}
	html, err := s.RenderToString(sReq.Template, &data, true, true)
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
	return s.Send(message)
}

func (s *MailerService) Send(message provider.Message) error {
	for _, p := range s.config.Providers {
		color.Cyan("Sending mail using %s provider...\n", p.Name())
		err := p.Send(message)
		if err == nil {
			return nil
		}
		color.Red("Failed to send message with %s provider. Attempting next provider...", p.Name())
	}
	return nil
}
