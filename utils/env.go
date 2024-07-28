// utils/env.go
package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvEngine struct{}

func NewEnvEngine() *EnvEngine {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &EnvEngine{}
}

func (*EnvEngine) LoadEnv(key string) string {
	return os.Getenv(key)
}
