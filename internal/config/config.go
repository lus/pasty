package config

import (
	"strings"
	"time"

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
	env.Load()

	Current = &Config{
		WebAddress:          env.MustString("WEB_ADDRESS", ":8080"),
		StorageType:         shared.StorageType(strings.ToLower(env.MustString("STORAGE_TYPE", "file"))),
		HastebinSupport:     env.MustBool("HASTEBIN_SUPPORT", false),
		IDLength:            env.MustInt("ID_LENGTH", 6),
		DeletionTokenLength: env.MustInt("DELETION_TOKEN_LENGTH", 12),
		RateLimit:           env.MustString("RATE_LIMIT", "30-M"),
		AutoDelete: &AutoDeleteConfig{
			Enabled:      env.MustBool("AUTODELETE", false),
			Lifetime:     env.MustDuration("AUTODELETE_LIFETIME", 720*time.Hour),
			TaskInterval: env.MustDuration("AUTODELETE_TASK_INTERVAL", 5*time.Minute),
		},
		File: &FileConfig{
			Path: env.MustString("STORAGE_FILE_PATH", "./data"),
		},
		Postgres: &PostgresConfig{
			DSN: env.MustString("STORAGE_POSTGRES_DSN", "postgres://pasty:pasty@localhost/pasty"),
		},
		MongoDB: &MongoDBConfig{
			DSN:        env.MustString("STORAGE_MONGODB_CONNECTION_STRING", "mongodb://pasty:pasty@localhost/pasty"),
			Database:   env.MustString("STORAGE_MONGODB_DATABASE", "pasty"),
			Collection: env.MustString("STORAGE_MONGODB_COLLECTION", "pastes"),
		},
		S3: &S3Config{
			Endpoint:        env.MustString("STORAGE_S3_ENDPOINT", ""),
			AccessKeyID:     env.MustString("STORAGE_S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: env.MustString("STORAGE_S3_SECRET_ACCESS_KEY", ""),
			SecretToken:     env.MustString("STORAGE_S3_SECRET_TOKEN", ""),
			Secure:          env.MustBool("STORAGE_S3_SECURE", true),
			Region:          env.MustString("STORAGE_S3_REGION", ""),
			Bucket:          env.MustString("STORAGE_S3_BUCKET", "pasty"),
		},
	}
}
