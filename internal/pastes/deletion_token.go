package pastes

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"math/rand"
	"strconv"
)

// deletionTokenContents represents the characters a deletion token may contain
const deletionTokenContents = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#+-.,"

// generateDeletionToken generates a new deletion token
func generateDeletionToken() (string, error) {
	// Read the deletion token length
	rawLength := env.Get("DELETION_TOKEN_LENGTH", "12")
	length, err := strconv.Atoi(rawLength)
	if err != nil {
		return "", err
	}

	// Generate the deletion token
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = deletionTokenContents[rand.Int63()%int64(len(deletionTokenContents))]
	}
	return string(bytes), nil
}
