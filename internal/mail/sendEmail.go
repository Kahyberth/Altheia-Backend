package mail

import (
	"fmt"
	"log"
	"net/smtp"
)

// SMTPClient defines the interface for sending emails
type SMTPClient interface {
	SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

// DefaultSMTPClient implements SMTPClient using the standard smtp package
type DefaultSMTPClient struct{}

func (c *DefaultSMTPClient) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, auth, from, to, msg)
}

// Global SMTP client instance
var smtpClient SMTPClient = &DefaultSMTPClient{}

// SetSMTPClient allows setting a custom SMTP client (used in tests)
func SetSMTPClient(client SMTPClient) {
	smtpClient = client
}

func SendWelcomeMessage(username string, to []string) error {
	if username == "" || len(to) == 0 || to[0] == "" {
		return fmt.Errorf("invalid parameters for sending email")
	}

	go func() {
		auth, emailUsername, emailHost := EmailConfig()

		subject := fmt.Sprintf("👋 ¡Bienvenido/a %s! Tu cuenta en Altheia EHR está lista", username)
		message := EmailTemplate(username)

		msg := "From: " + emailUsername + "\r\n" +
			"To: " + to[0] + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" + message

		err := smtpClient.SendMail(
			emailHost+":587",
			auth,
			emailUsername,
			to,
			[]byte(msg),
		)

		if err != nil {
			log.Printf("Error al enviar correo de bienvenida: %v", err)
		}
	}()

	return nil
}
