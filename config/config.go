package config

import (
	"os"
)

type Config struct {
	JWTSecret []byte
	Port      string
}

// LoadConfig reads variables from environment or falls back to defaults
func LoadConfig() *Config {
	jwtSecret := "my_super_secret_key"

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		JWTSecret: []byte(jwtSecret),
		Port:      port,
	}
}
