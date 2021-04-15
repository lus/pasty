package storage

import (
	"errors"
	"fmt"

	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/shared"
)

// Current holds the current storage driver
var Current Driver

// Driver represents a storage driver
type Driver interface {
	Initialize() error
	Terminate() error
	ListIDs() ([]string, error)
	Get(id string) (*shared.Paste, error)
	Save(paste *shared.Paste) error
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
func GetDriver(storageType shared.StorageType) (Driver, error) {
	switch storageType {
	case shared.StorageTypeFile:
		return new(FileDriver), nil
	case shared.StorageTypePostgres:
		// TODO: Implement Postgres driver
		return nil, errors.New("TODO")
	case shared.StorageTypeMongoDB:
		return new(MongoDBDriver), nil
	case shared.StorageTypeS3:
		return new(S3Driver), nil
	default:
		return nil, fmt.Errorf("invalid storage type '%s'", storageType)
	}
}
