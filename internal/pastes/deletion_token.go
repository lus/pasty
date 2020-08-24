package pastes

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/utils"
	"strconv"
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
