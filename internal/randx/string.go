package randx

import (
	"math/rand"
)

// String generates a random string with the given length.
func String(characters string, length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[rand.Int63()%int64(len(characters))]
	}
	return string(bytes)
}
