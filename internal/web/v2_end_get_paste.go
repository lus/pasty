package web

import (
	"errors"
	"github.com/lus/pasty/internal/pastes"
	"net/http"
)

func (server *Server) v2EndpointGetPaste(writer http.ResponseWriter, request *http.Request) {
	paste, ok := request.Context().Value("paste").(*pastes.Paste)
	if !ok {
		writeErr(request, writer, errors.New("missing paste object"))
		return
	}

	cpy := *paste
	cpy.ModificationToken = ""
	writeJSONOrErr(request, writer, http.StatusOK, cpy)
}
