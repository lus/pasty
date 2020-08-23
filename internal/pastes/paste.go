package pastes

import (
	"github.com/alexedwards/argon2id"
	"github.com/bwmarrin/snowflake"
)

func init() {
	snowflakeNode, _ = snowflake.NewNode(1)
}

// snowflakeNode holds the current snowflake node
var snowflakeNode *snowflake.Node

// Paste represents a saved paste
type Paste struct {
	ID                  snowflake.ID `json:"id" bson:"_id"`
	Content             string       `json:"content" bson:"content"`
	SuggestedSyntaxType string       `json:"suggestedSyntaxType" bson:"suggestedSyntaxType"`
	DeletionToken       string       `json:"deletionToken" bson:"deletionToken"`
}

// Create creates a new paste object using the given content
func Create(content string) (*Paste, error) {
	// TODO: Generate the suggested syntax type
	suggestedSyntaxType := ""

	// Generate the deletion token
	deletionToken, err := generateDeletionToken()
	if err != nil {
		return nil, err
	}

	// Return the paste object
	return &Paste{
		ID:                  snowflakeNode.Generate(),
		Content:             content,
		SuggestedSyntaxType: suggestedSyntaxType,
		DeletionToken:       deletionToken,
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
