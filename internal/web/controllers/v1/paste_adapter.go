package v1

import "github.com/lus/pasty/internal/shared"

type legacyPaste struct {
	ID            string `json:"id"`
	Content       string `json:"content"`
	DeletionToken string `json:"deletionToken,omitempty"`
	Created       int64  `json:"created"`
	AutoDelete    bool   `json:"autoDelete"`
}

func legacyFromModern(paste *shared.Paste) *legacyPaste {
	deletionToken := paste.ModificationToken
	if deletionToken == "" {
		deletionToken = paste.DeletionToken
	}

	return &legacyPaste{
		ID:            paste.ID,
		Content:       paste.Content,
		DeletionToken: deletionToken,
		Created:       paste.Created,
		AutoDelete:    paste.AutoDelete,
	}
}
