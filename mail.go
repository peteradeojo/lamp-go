package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailMesage struct {
	To        []string      `json:"recipients"`
	Subject   string        `json:"subject"`
	Message   template.HTML `json:"message"`
	PlainText string        `json:"plain"`
}

func SetupMailConfig() {
	smtpConfig.Host = os.Getenv("SMTP_HOST")
	smtpConfig.Port, _ = strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 16)
	smtpConfig.Username = os.Getenv("SMTP_USERNAME")
	smtpConfig.Password = os.Getenv("SMTP_PASSWORD")
}

func renderTemplate(htmlMessage MailMesage) (string, error) {
	tmpl, err := template.ParseFiles("email_template.html")
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, htmlMessage); err != nil {
		return "", err
	}

	return rendered.String(), nil
}

func SendMail(message MailMesage) error {
	renderedHTML, err := renderTemplate(message)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "Doppler <"+smtpConfig.Username+">")
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/plain", message.PlainText)
	m.AddAlternative("text/html", renderedHTML)

	for _, e := range message.To {
		// fmt.Println(i, e)
		m.SetHeader("To", e)

		d := gomail.NewDialer(smtpConfig.Host, int(smtpConfig.Port), smtpConfig.Username, smtpConfig.Password)
		err = d.DialAndSend(m)

		if err != nil {
			continue
		}
	}

	return nil
}
