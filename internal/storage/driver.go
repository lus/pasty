package storage

import (
	"fmt"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"strings"
)

// Current holds the current storage driver
var Current Driver

// Driver represents a storage driver
type Driver interface {
	Initialize() error
	Terminate() error
	ListIDs() ([]string, error)
	Get(id string) (*pastes.Paste, error)
	Save(paste *pastes.Paste) error
	Delete(id string) error
	Cleanup() (int, error)
}

// Load loads the current storage driver
func Load() error {
	// Define the driver to use
	storageType := env.Get("STORAGE_TYPE", "file")
	driver, err := GetDriver(storageType)
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

// GetDriver returns the driver with the given type string if it exists
func GetDriver(storageType string) (Driver, error) {
	switch strings.ToLower(storageType) {
	case "file":
		return new(FileDriver), nil
	case "s3":
		return new(S3Driver), nil
	case "mongodb":
		return new(MongoDBDriver), nil
	case "sql":
		return new(SQLDriver), nil
	default:
		return nil, fmt.Errorf("invalid storage type '%s'", storageType)
	}
}
