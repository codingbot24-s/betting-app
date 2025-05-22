package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


func GetEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error getting environment variable %s", key)
	}
	return value
}