package service

import (
	"github.com/carlosstrand/mailer/provider"
	"github.com/carlosstrand/mailer/types"
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"strings"
	"sync"
)

func (m *MailerService) SendFromReq(sReq types.SendMailFromTemplateRequest) error {
	data := sync.Map{}
	for key, value := range sReq.Vars {
		data.Store(key, value)
	}

	// Render Subject
	d := data
	d.Store("mode", "subject")
	subject, err := m.RenderToString(sReq.Template, &d, false, false)
	subject = strings.Trim(subject, "\t \n")

	// Render PlainText
	d = data
	d.Store("mode", "plain_text")
	plainText, err := m.RenderToString(sReq.Template, &d, false, false)
	if err != nil {
		return err
	}
	plainText = strings.Trim(plainText, "\t \n")

	// Render HTML
	d = data
	d.Store("mode", "html")
	html, err := m.RenderToString(sReq.Template, &d, true, true)
	if err != nil {
		return err
	}

	message := provider.Message{
		From:      sReq.From,
		To:        sReq.To,
		Subject:   subject,
		PlainText: plainText,
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
		color.Red("Send Error: %s", err.Error())
		color.Red("Failed to send message with %s provider. Attempting next provider...", p.Name())
	}
	return nil
}
