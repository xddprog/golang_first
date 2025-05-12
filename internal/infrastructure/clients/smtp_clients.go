package clients

import (
	"bytes"
	"golang/internal/infrastructure/config"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)


type SmtpClient struct {
	Dialer *gomail.Dialer
}


func NewSmtpClient() *SmtpClient {
	cfg := config.LoadSmtpConfig()
	
	return &SmtpClient{Dialer: gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)}
}


func (client *SmtpClient) SendInviteToDocument(to string, subject string, code string, documentTitle string, documentId string) error {
	message := gomail.NewMessage()
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	
	htmlBytes, err := os.ReadFile("invite_member.html")
	if err != nil {
		return err
	}

	template, err := template.New("invite").Parse(string(htmlBytes))
	if err != nil {
		return err
	}

	data := map[string]string{
		"DocumentTitle": 	documentTitle,
		"AccessCode":       code,
		"DocumentId": 		documentId,
	}

	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return err
	}
	message.SetBody("text/html", buf.String())

	if err := client.Dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}