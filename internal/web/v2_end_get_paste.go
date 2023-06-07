package web

import (
	"github.com/lus/pasty/internal/pastes"
	"net/http"
)

func (server *Server) v2EndpointGetPaste(writer http.ResponseWriter, request *http.Request) {
	paste, ok := request.Context().Value("paste").(*pastes.Paste)
	if !ok {
		writeString(writer, http.StatusInternalServerError, "missing paste object")
		return
	}

	cpy := *paste
	cpy.ModificationToken = ""
	writeJSONOrErr(writer, http.StatusOK, cpy)
}
