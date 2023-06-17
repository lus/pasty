package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	LogLevel                  string            `default:"info" split_words:"true"`
	Address                   string            `default:":8080" split_words:"true"`
	StorageDriver             string            `default:"sqlite" split_words:"true"`
	PasteIDLength             int               `default:"6" split_words:"true"`
	PasteIDCharset            string            `default:"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" split_words:"true"`
	ModificationTokensEnabled bool              `default:"true" split_words:"true"`
	ModificationTokenMaster   string            `split_words:"true"`
	ModificationTokenLength   int               `default:"12" split_words:"true"`
	ModificationTokenCharset  string            `default:"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" split_words:"true"`
	RateLimit                 string            `default:"30-M" split_words:"true"`
	PasteLengthCap            int               `default:"50000" split_words:"true"`
	AutoDelete                *AutoDeleteConfig `split_words:"true"`
	Reports                   *ReportConfig
	Postgres                  *PostgresConfig
	SQLite                    *SQLiteConfig
}

type AutoDeleteConfig struct {
	Enabled      bool          `default:"false"`
	Lifetime     time.Duration `default:"720h"`
	TaskInterval time.Duration `default:"5m" split_words:"true"`
}

type ReportConfig struct {
	Enabled      bool   `default:"false" split_words:"true"`
	WebhookURL   string `split_words:"true"`
	WebhookToken string `split_words:"true"`
}

type PostgresConfig struct {
	DSN string `default:"postgres://pasty:pasty@localhost/pasty"`
}

type SQLiteConfig struct {
	File string `default:":memory:"`
}

func Load() (*Config, error) {
	_ = godotenv.Overload()
	cfg := new(Config)
	if err := envconfig.Process("pasty", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
