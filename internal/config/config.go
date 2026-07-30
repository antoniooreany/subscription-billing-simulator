package config

import (
	"os"
	"strings"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	AutoMigrate bool
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	autoMigrate := strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE")))
	return Config{
		AppPort:     port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AutoMigrate: autoMigrate == "" || autoMigrate == "true" || autoMigrate == "1" || autoMigrate == "yes",
	}
}