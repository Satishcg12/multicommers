package email

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

type (
	EmailMessage struct {
		From    string
		To      string
		Subject string
		Body    string
	}
	EmailDaemon struct {
		dialer *gomail.Dialer
		msgCh  chan EmailMessage
	}
	EmailDaemonInterface interface {
		Start()
		Send(EmailMessage)
		SendEmail(EmailMessage)
	}
)

func NewEmailDaemon(Host, Port, Username, Password string) EmailDaemonInterface {
	port, _ := strconv.Atoi(Port)
	dialer := gomail.NewDialer(Host, port, Username, Password)
	return &EmailDaemon{
		dialer: dialer,
		msgCh:  make(chan EmailMessage, 100),
	}

}

func (e *EmailDaemon) Start() {
	go func() {
		for email := range e.msgCh {
			e.SendEmail(email)
		}
	}()
}

func (e *EmailDaemon) Send(email EmailMessage) {
	e.msgCh <- email
}

func (e *EmailDaemon) SendEmail(email EmailMessage) {
	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	if err := e.dialer.DialAndSend(m); err != nil {
		panic(err)
	}
}
