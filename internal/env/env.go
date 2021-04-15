package env

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/lus/pasty/internal/static"
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

// Bool uses Get and parses it into a boolean
func Bool(key string, fallback bool) bool {
	parsed, _ := strconv.ParseBool(Get(key, strconv.FormatBool(fallback)))
	return parsed
}

// Duration uses Get and parses it into a duration
func Duration(key string, fallback time.Duration) time.Duration {
	parsed, _ := time.ParseDuration(Get(key, fallback.String()))
	return parsed
}
