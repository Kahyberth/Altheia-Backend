package mail

import (
	"Altheia-Backend/config"
	"net/smtp"
	"os"
	"sync"
)

var (
	once          sync.Once
	authEmail     smtp.Auth
	emailHost     string
	emailUsername string
	emailPassword string
)

func EmailConfig() (smtp.Auth, string, string) {
	config.LoadEnv()
	once.Do(func() {
		emailPassword = os.Getenv("EMAIL_PASSWORD")
		emailUsername = os.Getenv("EMAIL_USERNAME")
		emailHost = os.Getenv("EMAIL_HOST")
		authEmail = smtp.PlainAuth("", emailUsername, emailPassword, emailHost)
	})
	return authEmail, emailUsername, emailHost
}
