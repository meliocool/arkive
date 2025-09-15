package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBUser, DBPassword, DBName, DBHost                          string
	ZohoUser, ZohoPassword, ZohoHost, ZohoServiceName, ZohoPort string
	JwtSecret                                                   string
	IPFSAPIKey, IPFSAPISecret                                   string
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
		return nil, fmt.Errorf("missing one or more required DB env vars")
	}

	ZohoUser := os.Getenv("EMAIL_SMTP_USER")
	ZohoPassword := os.Getenv("EMAIL_SMTP_PASS")
	ZohoHost := os.Getenv("EMAIL_SMTP_HOST")
	ZohoPort := os.Getenv("EMAIL_SMTP_PORT")

	if ZohoUser == "" || ZohoPassword == "" || ZohoHost == "" || ZohoPort == "" {
		return nil, fmt.Errorf("missing one or more required Email SMTP env vars")
	}

	JwtSecret := os.Getenv("JWT_SECRET")

	IPFSAPIKey := os.Getenv("IPFS_API_KEY")
	IPFSAPISecret := os.Getenv("IPFS_API_SECRET")

	if IPFSAPIKey == "" || IPFSAPISecret == "" {
		return nil, fmt.Errorf("missing one or more required IPFS env vars")
	}

	cfg := &Config{
		DBUser:     DBUser,
		DBPassword: DBPassword,
		DBName:     DBName,
		DBHost:     DBHost,

		ZohoUser:     ZohoUser,
		ZohoPassword: ZohoPassword,
		ZohoHost:     ZohoHost,
		ZohoPort:     ZohoPort,

		JwtSecret: JwtSecret,

		IPFSAPIKey:    IPFSAPIKey,
		IPFSAPISecret: IPFSAPISecret,
	}

	return cfg, nil
}
