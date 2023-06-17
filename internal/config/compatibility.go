package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

var removedKeys = []string{
	"PASTY_HASTEBIN_SUPPORT",
	"PASTY_STORAGE_FILE_PATH",
	"PASTY_STORAGE_MONGODB_CONNECTION_STRING",
	"PASTY_STORAGE_MONGODB_DATABASE",
	"PASTY_STORAGE_MONGODB_COLLECTION",
	"PASTY_STORAGE_S3_ENDPOINT",
	"PASTY_STORAGE_S3_ACCESS_KEY_ID",
	"PASTY_STORAGE_S3_SECRET_ACCESS_KEY",
	"PASTY_STORAGE_S3_SECRET_TOKEN",
	"PASTY_STORAGE_S3_SECURE",
	"PASTY_STORAGE_S3_REGION",
	"PASTY_STORAGE_S3_BUCKET",
}

var keyRedirects = map[string][]string{
	"PASTY_ADDRESS":                     {"PASTY_WEB_ADDRESS"},
	"PASTY_STORAGE_DRIVER":              {"PASTY_STORAGE_TYPE"},
	"PASTY_POSTGRES_DSN":                {"PASTY_STORAGE_POSTGRES_DSN"},
	"PASTY_PASTE_ID_LENGTH":             {"PASTY_ID_LENGTH"},
	"PASTY_PASTE_ID_CHARSET":            {"PASTY_ID_CHARACTERS"},
	"PASTY_PASTE_LENGTH_CAP":            {"PASTY_LENGTH_CAP"},
	"PASTY_REPORTS_ENABLED":             {"PASTY_REPORTS_ENABLED"},
	"PASTY_REPORTS_WEBHOOK_URL":         {"PASTY_REPORT_WEBHOOK"},
	"PASTY_REPORTS_WEBHOOK_TOKEN":       {"PASTY_REPORT_WEBHOOK_TOKEN"},
	"PASTY_CLEANUP_ENABLED":             {"PASTY_AUTODELETE"},
	"PASTY_CLEANUP_PASTE_LIFETIME":      {"PASTY_AUTODELETE_LIFETIME"},
	"PASTY_CLEANUP_TASK_INTERVAL":       {"PASTY_AUTODELETE_TASK_INTERVAL"},
	"PASTY_MODIFICATION_TOKENS_ENABLED": {"PASTY_MODIFICATION_TOKENS", "PASTY_DELETION_TOKENS"},
	"PASTY_MODIFICATION_TOKEN_CHARSET":  {"PASTY_MODIFICATION_TOKEN_CHARACTERS"},
	"PASTY_MODIFICATION_TOKEN_MASTER":   {"PASTY_DELETION_TOKEN_MASTER"},
	"PASTY_MODIFICATION_TOKEN_LENGTH":   {"PASTY_DELETION_TOKEN_LENGTH"},
}

// Compatibility runs several compatibility measurements.
// This is used to redirect legacy config keys to their new equivalent or print warnings about deprecated ones.
func Compatibility() {
	_ = godotenv.Overload()

	for _, key := range removedKeys {
		if isSet(key) {
			log.Warn().Msgf("You have set the '%s' environment variable. This variable has been discontinued and has no further effect.", key)
		}
	}

	for newKey, oldKeys := range keyRedirects {
		if !isSet(newKey) {
			for _, oldKey := range oldKeys {
				if isSet(oldKey) {
					if err := os.Setenv(newKey, os.Getenv(oldKey)); err != nil {
						continue
					}
					log.Warn().Msgf("You have set the '%s' environment variable. This variable has been renamed to '%s'. The value has been propagated, but please consider adjusting your configuration to avoid further complications.", oldKey, newKey)
					break
				}
			}
		}
	}
}

func isSet(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}
