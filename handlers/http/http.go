package http

import "github.com/carlosstrand/mailer/service"

type HTTPHandler struct {
	mailer *service.MailerService
}

func NewHTTPHandler(mailer *service.MailerService) *HTTPHandler {
	return &HTTPHandler{mailer: mailer}
}
