package web

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (server *Server) v2MiddlewareInjectPaste(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		pasteID := strings.TrimSpace(chi.URLParam(request, "paste_id"))
		if pasteID == "" {
			writeString(writer, http.StatusNotFound, "paste not found")
			return
		}

		paste, err := server.Storage.Pastes().FindByID(request.Context(), pasteID)
		if err != nil {
			writeErr(request, writer, err)
		}
		if paste == nil {
			writeString(writer, http.StatusNotFound, "paste not found")
			return
		}

		if paste.Metadata == nil {
			paste.Metadata = make(map[string]any)
		}

		request = request.WithContext(context.WithValue(request.Context(), "paste", paste))

		next.ServeHTTP(writer, request)
	})
}
