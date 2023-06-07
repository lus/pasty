package storage

import (
	"context"
	"github.com/lus/pasty/internal/config"
	"strings"

	"github.com/lus/pasty/internal/paste"
	"github.com/lus/pasty/internal/storage/postgres"
)

// Driver represents a storage driver
type Driver interface {
	Initialize(ctx context.Context, cfg *config.Config) error
	Close() error
	ListIDs() ([]string, error)
	Get(id string) (*paste.Paste, error)
	Save(paste *paste.Paste) error
	Delete(id string) error
	Cleanup() (int, error)
}

// ResolveDriver returns the driver with the given name if it exists
func ResolveDriver(name string) (Driver, bool) {
	switch strings.TrimSpace(strings.ToLower(name)) {
	case "postgres":
		return new(postgres.Driver), true
	default:
		return nil, false
	}
}
