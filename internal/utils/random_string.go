package utils

import (
	"math/rand"
	"time"
)

// stringContents holds the chars a random string can contain
const stringContents = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString returns a random string with the given length
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = stringContents[rand.Int63()%int64(len(stringContents))]
	}
	return string(bytes)
}
