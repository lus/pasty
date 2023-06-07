package pastes

import "github.com/alexedwards/argon2id"

type Paste struct {
	ID                string         `json:"id"`
	Content           string         `json:"content"`
	ModificationToken string         `json:"modificationToken,omitempty"`
	Created           int64          `json:"created"`
	Metadata          map[string]any `json:"metadata"`
}

func (paste *Paste) HashModificationToken() error {
	if paste.ModificationToken == "" {
		return nil
	}
	hash, err := argon2id.CreateHash(paste.ModificationToken, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	paste.ModificationToken = hash
	return nil
}

func (paste *Paste) CheckModificationToken(modificationToken string) bool {
	if paste.ModificationToken == "" {
		return false
	}
	match, err := argon2id.ComparePasswordAndHash(modificationToken, paste.ModificationToken)
	return err == nil && match
}
