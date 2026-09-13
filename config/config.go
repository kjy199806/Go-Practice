package config

import (
	"os"
)

type Config struct {
	JWTSecret []byte
	Port      string
}

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
