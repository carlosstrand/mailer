package sendgrid

import (
	"errors"
	"github.com/carlosstrand/mailer/provider"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

type SendgridProvider struct {
	client *sendgrid.Client
}

func NewSendgridProvider(token string) *SendgridProvider {
	client := sendgrid.NewSendClient(token)
	return &SendgridProvider{
		client: client,
	}
}

func (p *SendgridProvider) Name() string {
	return "SendGrid"
}

// Send a email message
func (p *SendgridProvider) Send(msg provider.Message) error {
	from := mail.NewEmail(msg.From.Name, msg.From.Email)
	to := mail.NewEmail(msg.To.Name, msg.To.Email)
	sgMessage := mail.NewSingleEmail(from, msg.Subject, to, msg.PlainText, msg.HTML)
	res, err := p.client.Send(sgMessage)
	if err != nil {
		return err
	}
	if res.StatusCode > 200 {
		return errors.New(http.StatusText(res.StatusCode))
	}
	return nil
}
