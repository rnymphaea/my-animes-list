package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get(key string) (string, bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ")
	}

	return os.LookupEnv(key)
}
