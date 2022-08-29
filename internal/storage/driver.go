package storage

import (
	"fmt"
	"strings"

	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/paste"
	"github.com/lus/pasty/internal/storage/file"
	"github.com/lus/pasty/internal/storage/mongodb"
	"github.com/lus/pasty/internal/storage/postgres"
	"github.com/lus/pasty/internal/storage/s3"
)

// Current holds the current storage driver
var Current Driver

// Driver represents a storage driver
type Driver interface {
	Initialize() error
	Terminate() error
	ListIDs() ([]string, error)
	Get(id string) (*paste.Paste, error)
	Save(paste *paste.Paste) error
	Delete(id string) error
	Cleanup() (int, error)
}

// Load loads the current storage driver
func Load() error {
	// Define the driver to use
	driver, err := GetDriver(config.Current.StorageType)
	if err != nil {
		return err
	}

	// Initialize the driver
	err = driver.Initialize()
	if err != nil {
		return err
	}
	Current = driver
	return nil
}

// GetDriver returns the driver with the given type if it exists
func GetDriver(storageType string) (Driver, error) {
	switch strings.TrimSpace(strings.ToLower(storageType)) {
	case "file":
		return new(file.FileDriver), nil
	case "postgres":
		return new(postgres.PostgresDriver), nil
	case "mongodb":
		return new(mongodb.MongoDBDriver), nil
	case "s3":
		return new(s3.S3Driver), nil
	default:
		return nil, fmt.Errorf("invalid storage type '%s'", storageType)
	}
}
