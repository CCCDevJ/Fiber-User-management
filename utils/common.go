package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
)

var (
	Store         *session.Store
	USERNAME      = "USERNAME"
	PASSWORD      = "PASSWORD"
	DATABASE_NAME = "DATABASE_NAME"
	HOST          = "HOST"
	PORT          = "PORT"
	APP_NAME      = "APP_NAME"
	WEB_PORT      = "WEB_PORT"
	BASE_URL      = "BASE_URL"
	EMAIL         = "email"
	ID            = "id"
	NAME          = "name"
	ROLE          = "role"
)

// use godot package to load/read the .env file and
// return the value of the key
func DotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetEncryptionSHA512(key string) string {
	h := sha512.New()
	h.Write([]byte(key))
	sha := h.Sum(nil) // "sha" is uint8 type, encoded in base16

	return hex.EncodeToString(sha)
}
