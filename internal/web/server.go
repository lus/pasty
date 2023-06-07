package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/lus/pasty/internal/storage"
	"net/http"
)

type Server struct {
	// The address the web server should listen to.
	Address string

	// The storage driver to use.
	Storage storage.Driver

	// Whether the Hastebin support should be enabled.
	// If this is set to 'false', the Hastebin specific endpoints will not be registered.
	HastebinSupport bool

	// The length of newly generated paste IDs.
	PasteIDLength int
	// The charset to use when generating new paste IDs.
	PasteIDCharset string

	// The maximum length of newly generated pastes.
	PasteLengthCap int

	// Whether modification tokens are enabled.
	ModificationTokensEnabled bool
	// The length of newly generated modification tokens.
	ModificationTokenLength int
	// The charset to use when generating new modification tokens.
	ModificationTokenCharset string

	// The administration tokens.
	AdminTokens []string
}

func (server *Server) Start() error {
	router := chi.NewRouter()

	// Register the paste API endpoints
	router.With(server.v2MiddlewareInjectPaste).Get("/api/v2/pastes/{paste_id}", server.v2EndpointGetPaste)
	router.Post("/api/v2/pastes", server.v2EndpointCreatePaste)
	router.With(server.v2MiddlewareInjectPaste, server.v2MiddlewareAuthorize).Patch("/api/v2/pastes/{paste_id}", server.v2EndpointModifyPaste)
	router.Delete("/api/v2/pastes/{paste_id}", server.v2EndpointDeletePaste)

	return http.ListenAndServe(server.Address, router)
}
