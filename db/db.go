package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				break
			}
		}

		if attempts < maxAttempts {
			log.Printf("Attempt %d/%d: Database not ready, retrying in 5 seconds...", attempts, maxAttempts)
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxAttempts, err)
	}

	log.Println("Successfully connected to database")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
