package storage

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/utils"
	"strconv"
)

// AcquireID generates a new unique ID
func AcquireID() (string, error) {
	// Read the ID length
	rawLength := env.Get("ID_LENGTH", "6")
	length, err := strconv.Atoi(rawLength)
	if err != nil {
		return "", err
	}

	// Generate the unique ID
	for {
		id := utils.RandomString(length)
		paste, err := Current.Get(id)
		if err != nil {
			return "", err
		}
		if paste == nil {
			return id, nil
		}
	}
}
