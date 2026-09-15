package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Load reads a local .env if present. Docker and CI already inject env vars,
// so a missing file is not an error.
func Load() {
	_ = godotenv.Load()
}

func Get(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
