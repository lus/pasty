package utils

import (
	"math/rand"
	"time"
)

// RandomString returns a random string with the given length
func RandomString(characters string, length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = characters[rand.Int63()%int64(len(characters))]
	}
	return string(bytes)
}
