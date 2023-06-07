package pastes

type Paste struct {
	ID                string                 `json:"id"`
	Content           string                 `json:"content"`
	ModificationToken string                 `json:"modificationToken,omitempty"`
	Created           int64                  `json:"created"`
	Metadata          map[string]interface{} `json:"metadata"`
}
