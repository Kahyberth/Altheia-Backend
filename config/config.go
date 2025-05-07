package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
