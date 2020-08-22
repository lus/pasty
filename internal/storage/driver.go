package storage

import (
	"fmt"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/bwmarrin/snowflake"
	"strings"
)

// Current holds the current storage driver
var Current Driver

// Driver represents a storage driver
type Driver interface {
	Initialize() error
	Terminate() error
	Get(id snowflake.ID) (*pastes.Paste, error)
	Save(paste *pastes.Paste) error
	Delete(id snowflake.ID) error
}

// Load loads the current storage driver
func Load() error {
	storageType := strings.ToLower(env.Get("STORAGE_TYPE", "file"))
	switch storageType {
	case "file":
		driver := new(FileDriver)
		err := driver.Initialize()
		if err != nil {
			return err
		}
		Current = driver
		return nil
	default:
		return fmt.Errorf("invalid storage type '%s'", storageType)
	}
}
