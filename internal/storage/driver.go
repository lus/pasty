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
}

// Load loads the current storage driver
func Load() error {
	// Define the driver to use
	var driver Driver
	storageType := strings.ToLower(env.Get("STORAGE_TYPE", "file"))
	switch storageType {
	case "file":
		driver = new(FileDriver)
		break
	case "s3":
		driver = new(S3Driver)
		break
	case "mongodb":
		driver = new(MongoDBDriver)
		break
	default:
		return fmt.Errorf("invalid storage type '%s'", storageType)
	}

	// Initialize the driver
	err := driver.Initialize()
	if err != nil {
		return err
	}
	Current = driver
	return nil
}
