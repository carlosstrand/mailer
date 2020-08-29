package provider

type EmailAccount struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Message struct {
	From      EmailAccount `json:"from"`
	To        EmailAccount `json:"to"`
	Subject   string       `json:"subject"`
	PlainText string       `json:"plain_text"`
	HTML      string       `json:"html"`
}

type MailProvider interface {
	Name() string
	Send(Message) error
}
