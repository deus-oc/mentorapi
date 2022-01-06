package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(key string) string {

	// load .env file
	err := godotenv.Load("config.env")

	if err != nil {
		log.Println(err)
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
