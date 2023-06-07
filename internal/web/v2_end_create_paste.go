package web

import (
	"encoding/json"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/randx"
	"io"
	"net/http"
	"time"
)

type v2EndpointCreatePastePayload struct {
	Content  string         `json:"content"`
	Metadata map[string]any `json:"metadata"`
}

func (server *Server) v2EndpointCreatePaste(writer http.ResponseWriter, request *http.Request) {
	// Read, parse and validate the request payload
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeErr(writer, err)
		return
	}
	payload := new(v2EndpointCreatePastePayload)
	if err := json.Unmarshal(body, payload); err != nil {
		writeErr(writer, err)
		return
	}
	if payload.Content == "" {
		writeString(writer, http.StatusBadRequest, "missing paste content")
		return
	}
	if server.PasteLengthCap > 0 && len(payload.Content) > server.PasteLengthCap {
		writeString(writer, http.StatusBadRequest, "too large paste content")
		return
	}

	id, err := pastes.GenerateID(request.Context(), server.Storage.Pastes(), server.PasteIDCharset, server.PasteIDLength)
	if err != nil {
		writeErr(writer, err)
		return
	}

	paste := &pastes.Paste{
		ID:       id,
		Content:  payload.Content,
		Created:  time.Now().Unix(),
		Metadata: payload.Metadata,
	}

	modificationToken := ""
	if server.ModificationTokensEnabled {
		modificationToken = randx.String(server.ModificationTokenCharset, server.ModificationTokenLength)
		paste.ModificationToken = modificationToken

		if err := paste.HashModificationToken(); err != nil {
			writeErr(writer, err)
			return
		}
	}

	if err := server.Storage.Pastes().Upsert(request.Context(), paste); err != nil {
		writeErr(writer, err)
		return
	}

	cpy := *paste
	cpy.ModificationToken = modificationToken
	writeJSONOrErr(writer, http.StatusCreated, cpy)
}
