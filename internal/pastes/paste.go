package pastes

import "github.com/bwmarrin/snowflake"

// Paste represents a saved paste
type Paste struct {
	ID                  snowflake.ID `json:"id" bson:"_id"`
	Content             string       `json:"content" bson:"content"`
	SuggestedSyntaxType string       `json:"suggestedSyntaxType" bson:"suggestedSyntaxType"`
	DeletionToken       string       `json:"deletionToken" bson:"deletionToken"`
}
