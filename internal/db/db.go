package db

import (
	"Altheia-Backend/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *gorm.DB
)

func GetDB() *gorm.DB {
	config.LoadEnv()
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to DB: %v", err)
		}
		instance = db
		fmt.Println("Connected to DB")
	})

	return instance
}
