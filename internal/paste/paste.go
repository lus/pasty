package paste

import (
	"github.com/alexedwards/argon2id"
)

// Paste represents a paste
type Paste struct {
	ID                string                 `json:"id"`
	Content           string                 `json:"content"`
	ModificationToken string                 `json:"modificationToken,omitempty"`
	Created           int64                  `json:"created"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// HashModificationToken hashes the current modification token of a paste
func (paste *Paste) HashModificationToken() error {
	hash, err := argon2id.CreateHash(paste.ModificationToken, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	paste.ModificationToken = hash
	return nil
}

// CheckModificationToken checks whether the given modification token is correct
func (paste *Paste) CheckModificationToken(modificationToken string) bool {
	match, err := argon2id.ComparePasswordAndHash(modificationToken, paste.ModificationToken)
	return err == nil && match
}
