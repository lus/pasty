package web

import (
	"github.com/lus/pasty/internal/pastes"
	"net/http"
)

func (server *Server) v2EndpointDeletePaste(writer http.ResponseWriter, request *http.Request) {
	paste, ok := request.Context().Value("paste").(*pastes.Paste)
	if !ok {
		writeString(writer, http.StatusInternalServerError, "missing paste object")
		return
	}

	if err := server.Storage.Pastes().DeleteByID(request.Context(), paste.ID); err != nil {
		writeErr(request, writer, err)
	}
	writer.WriteHeader(http.StatusOK)
}
