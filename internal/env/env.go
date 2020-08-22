package env

import (
	"github.com/Lukaesebrot/pasty/internal/static"
	"github.com/joho/godotenv"
	"os"
)

// Load loads an optional .env file
func Load() {
	godotenv.Load()
}

// Get returns the content of the environment variable with the given key or the given fallback
func Get(key, fallback string) string {
	found := os.Getenv(static.EnvironmentVariablePrefix + key)
	if found == "" {
		return fallback
	}
	return found
}
