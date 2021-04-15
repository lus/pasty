package shared

import (
	"github.com/alexedwards/argon2id"
)

// Paste represents a saved paste
type Paste struct {
	ID            string `json:"id" bson:"_id"`
	Content       string `json:"content" bson:"content"`
	DeletionToken string `json:"deletionToken,omitempty" bson:"deletionToken"`
	Created       int64  `json:"created" bson:"created"`
	AutoDelete    bool   `json:"autoDelete" bson:"autoDelete"`
}

// HashDeletionToken hashes the current deletion token of a paste
func (paste *Paste) HashDeletionToken() error {
	hash, err := argon2id.CreateHash(paste.DeletionToken, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	paste.DeletionToken = hash
	return nil
}

// CheckDeletionToken checks whether or not the given deletion token is correct
func (paste *Paste) CheckDeletionToken(deletionToken string) bool {
	match, err := argon2id.ComparePasswordAndHash(deletionToken, paste.DeletionToken)
	return err == nil && match
}
