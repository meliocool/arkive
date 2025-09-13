package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBUser, DBPassword, DBName, DBHost string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	DBUser := os.Getenv("POSTGRES_USER")
	DBPassword := os.Getenv("POSTGRES_PASSWORD")
	DBName := os.Getenv("POSTGRES_DB")
	DBHost := os.Getenv("POSTGRES_HOST")

	if DBUser == "" || DBPassword == "" || DBName == "" || DBHost == "" {
		return nil, fmt.Errorf("missing one or more required env vars")
	}

	cfg := &Config{
		DBUser:     DBUser,
		DBPassword: DBPassword,
		DBName:     DBName,
		DBHost:     DBHost,
	}

	return cfg, nil
}
