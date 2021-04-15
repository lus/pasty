package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lus/pasty/internal/env"
	"github.com/lus/pasty/internal/shared"
)

// Config represents the general application configuration structure
type Config struct {
	WebAddress          string
	StorageType         shared.StorageType
	HastebinSupport     bool
	IDLength            int
	DeletionTokenLength int
	RateLimit           string
	AutoDelete          *AutoDeleteConfig
	File                *FileConfig
	Postgres            *PostgresConfig
	MongoDB             *MongoDBConfig
	S3                  *S3Config
}

// AutoDeleteConfig represents the configuration specific for the AutoDelete behaviour
type AutoDeleteConfig struct {
	Enabled      bool
	Lifetime     time.Duration
	TaskInterval time.Duration
}

// FileConfig represents the configuration specific for the file storage driver
type FileConfig struct {
	Path string
}

// PostgresConfig represents the configuration specific for the Postgres storage driver
type PostgresConfig struct {
	DSN string
}

// MongoDBConfig represents the configuration specific for the MongoDB storage driver
type MongoDBConfig struct {
	DSN        string
	Database   string
	Collection string
}

// S3Config represents the configuration specific for the S3 storage driver
type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SecretToken     string
	Secure          bool
	Region          string
	Bucket          string
}

// Current holds the currently loaded config
var Current *Config

// Load loads the current config from environment variables and an optional .env file
func Load() {
	godotenv.Load()

	Current = &Config{
		WebAddress:          env.MustString("PASTY_WEB_ADDRESS", ":8080"),
		StorageType:         shared.StorageType(strings.ToLower(env.MustString("PASTY_STORAGE_TYPE", "file"))),
		HastebinSupport:     env.MustBool("PASTY_HASTEBIN_SUPPORT", false),
		IDLength:            env.MustInt("PASTY_ID_LENGTH", 6),
		DeletionTokenLength: env.MustInt("PASTY_DELETION_TOKEN_LENGTH", 12),
		RateLimit:           env.MustString("PASTY_RATE_LIMIT", "30-M"),
		AutoDelete: &AutoDeleteConfig{
			Enabled:      env.MustBool("PASTY_AUTODELETE", false),
			Lifetime:     env.MustDuration("PASTY_AUTODELETE_LIFETIME", 720*time.Hour),
			TaskInterval: env.MustDuration("PASTY_AUTODELETE_TASK_INTERVAL", 5*time.Minute),
		},
		File: &FileConfig{
			Path: env.MustString("PASTY_STORAGE_FILE_PATH", "./data"),
		},
		Postgres: &PostgresConfig{
			DSN: env.MustString("PASTY_STORAGE_POSTGRES_DSN", "postgres://pasty:pasty@localhost/pasty"),
		},
		MongoDB: &MongoDBConfig{
			DSN:        env.MustString("PASTY_STORAGE_MONGODB_CONNECTION_STRING", "mongodb://pasty:pasty@localhost/pasty"),
			Database:   env.MustString("PASTY_STORAGE_MONGODB_DATABASE", "pasty"),
			Collection: env.MustString("PASTY_STORAGE_MONGODB_COLLECTION", "pastes"),
		},
		S3: &S3Config{
			Endpoint:        env.MustString("PASTY_STORAGE_S3_ENDPOINT", ""),
			AccessKeyID:     env.MustString("PASTY_STORAGE_S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: env.MustString("PASTY_STORAGE_S3_SECRET_ACCESS_KEY", ""),
			SecretToken:     env.MustString("PASTY_STORAGE_S3_SECRET_TOKEN", ""),
			Secure:          env.MustBool("PASTY_STORAGE_S3_SECURE", true),
			Region:          env.MustString("PASTY_STORAGE_S3_REGION", ""),
			Bucket:          env.MustString("PASTY_STORAGE_S3_BUCKET", "pasty"),
		},
	}
}
