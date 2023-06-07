package web

import (
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/slices"
	"net/http"
	"strings"
)

func (server *Server) v2MiddlewareAuthorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		paste, ok := request.Context().Value("paste").(*pastes.Paste)
		if !ok {
			writeString(writer, http.StatusInternalServerError, "missing paste object")
			return
		}

		authHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)
		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			writeString(writer, http.StatusUnauthorized, "unauthorized")
			return
		}

		isAdmin := slices.Contains(server.AdminTokens, authHeader[1])
		if isAdmin {
			next.ServeHTTP(writer, request)
			return
		}

		if !server.ModificationTokensEnabled || !paste.CheckModificationToken(authHeader[1]) {
			writeString(writer, http.StatusUnauthorized, "unauthorized")
			return
		}

		next.ServeHTTP(writer, request)
	})
}
