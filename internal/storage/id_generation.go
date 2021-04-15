package storage

import (
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/utils"
)

// AcquireID generates a new unique ID
func AcquireID() (string, error) {
	for {
		id := utils.RandomString(config.Current.IDLength)
		paste, err := Current.Get(id)
		if err != nil {
			return "", err
		}
		if paste == nil {
			return id, nil
		}
	}
}
