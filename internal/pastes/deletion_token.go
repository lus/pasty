package pastes

import (
	"strconv"

	"github.com/lus/pasty/internal/env"
	"github.com/lus/pasty/internal/utils"
)

// generateDeletionToken generates a new deletion token
func generateDeletionToken() (string, error) {
	// Read the deletion token length
	rawLength := env.Get("DELETION_TOKEN_LENGTH", "12")
	length, err := strconv.Atoi(rawLength)
	if err != nil {
		return "", err
	}

	// Generate the deletion token
	return utils.RandomString(length), nil
}
