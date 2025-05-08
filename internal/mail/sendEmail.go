package mail

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendWelcomeMessage(username string, to []string) error {

	if username == "" || len(to) == 0 || to[0] == "" {
		return fmt.Errorf("parámetros inválidos para enviar el correo")
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

		err := smtp.SendMail(
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
