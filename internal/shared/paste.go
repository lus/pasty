package shared

import (
	"log"

	"github.com/alexedwards/argon2id"
)

// Paste represents a saved paste
type Paste struct {
	ID                string                 `json:"id" bson:"_id"`
	Content           string                 `json:"content" bson:"content"`
	DeletionToken     string                 `json:"deletionToken,omitempty" bson:"deletionToken"` // Required for legacy paste storage support
	ModificationToken string                 `json:"modificationToken,omitempty" bson:"modificationToken"`
	Created           int64                  `json:"created" bson:"created"`
	Metadata          map[string]interface{} `json:"metadata" bson:"metadata"`
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

// CheckModificationToken checks whether or not the given modification token is correct
func (paste *Paste) CheckModificationToken(modificationToken string) bool {
	// The modification token may be stored in the deletion token field in old pastes
	usedToken := paste.ModificationToken
	if usedToken == "" {
		usedToken = paste.DeletionToken
		if usedToken != "" {
			log.Println("WARNING: You seem to have pastes with the old 'deletionToken' field stored in your storage driver. Though this does not cause any issues right now, it may in the future. Consider some kind of migration.")
		}
	}

	match, err := argon2id.ComparePasswordAndHash(modificationToken, usedToken)
	return err == nil && match
}
