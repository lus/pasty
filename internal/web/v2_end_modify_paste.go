package web

import (
	"encoding/json"
	"fmt"
	"github.com/lus/pasty/internal/maps"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/static"
	"io"
	"net/http"
)

type v2EndpointModifyPastePayload struct {
	Content  *string        `json:"content"`
	Metadata map[string]any `json:"metadata"`
}

func (server *Server) v2EndpointModifyPaste(writer http.ResponseWriter, request *http.Request) {
	paste, ok := request.Context().Value("paste").(*pastes.Paste)
	if !ok {
		writeString(writer, http.StatusInternalServerError, "missing paste object")
		return
	}

	// Read, parse and validate the request payload
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeErr(request, writer, err)
		return
	}
	payload := new(v2EndpointModifyPastePayload)
	if err := json.Unmarshal(body, payload); err != nil {
		writeErr(request, writer, err)
		return
	}
	if payload.Content != nil && *payload.Content == "" {
		writeString(writer, http.StatusBadRequest, "missing paste content")
		return
	}
	if payload.Content != nil && server.PasteLengthCap > 0 && len(*payload.Content) > server.PasteLengthCap {
		writeString(writer, http.StatusBadRequest, "too large paste content")
		return
	}
	if payload.Metadata != nil && maps.ExceedsDimensions(payload.Metadata, static.MaxMetadataWidth, static.MaxMetadataDepth) {
		writeString(writer, http.StatusBadRequest, fmt.Sprintf("metadata exceeds maximum dimensions (max. width: %d; max. depth: %d)", static.MaxMetadataWidth, static.MaxMetadataDepth))
		return
	}

	// Modify the paste itself
	if payload.Content != nil {
		paste.Content = *payload.Content
	}
	if payload.Metadata != nil {
		for key, value := range payload.Metadata {
			if value == nil {
				delete(paste.Metadata, key)
				continue
			}
			paste.Metadata[key] = value
		}
	}

	// Save the modified paste
	if err := server.Storage.Pastes().Upsert(request.Context(), paste); err != nil {
		writeErr(request, writer, err)
	}
}
