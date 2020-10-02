package pastes

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/alexedwards/argon2id"
	"time"
)

// Paste represents a saved paste
type Paste struct {
	ID                  string `json:"id" bson:"_id"`
	Content             string `json:"content" bson:"content"`
	SuggestedSyntaxType string `json:"suggestedSyntaxType" bson:"suggestedSyntaxType"`
	DeletionToken       string `json:"deletionToken" bson:"deletionToken,omitempty"`
	Created             int64  `json:"created" bson:"created"`
	AutoDelete          bool   `json:"autoDelete" bson:"autoDelete"`
}

// Create creates a new paste object using the given content
func Create(id, content string) (*Paste, error) {
	// TODO: Generate the suggested syntax type
	suggestedSyntaxType := ""

	// Generate the deletion token
	deletionToken, err := generateDeletionToken()
	if err != nil {
		return nil, err
	}

	// Return the paste object
	return &Paste{
		ID:                  id,
		Content:             content,
		SuggestedSyntaxType: suggestedSyntaxType,
		DeletionToken:       deletionToken,
		Created:             time.Now().Unix(),
		AutoDelete:          env.Bool("AUTODELETE", false),
	}, nil
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
