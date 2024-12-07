package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		panic(fmt.Sprintf("Error loading .env file: %v\n", err))
	}
	return os.Getenv(key)
}