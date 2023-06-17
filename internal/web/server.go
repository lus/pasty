package web

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/lus/pasty/internal/meta"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/reports"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/pkg/chiimplicitok"
	"github.com/lus/pasty/pkg/chizerolog"
	"net/http"
)

type Server struct {
	// The address the web server should listen to.
	Address string

	// The storage driver to use.
	Storage storage.Driver

	// The report client to use to send reports.
	// If this is set to nil, the report system will be considered deactivated.
	ReportClient *reports.Client

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

	httpServer *http.Server
}

func (server *Server) Start() error {
	router := chi.NewRouter()

	router.Use(chizerolog.Logger)
	router.Use(chizerolog.Recover)
	router.Use(chiimplicitok.Middleware)

	// Register the web frontend handler
	router.Get("/*", frontendHandler(router.NotFoundHandler()))

	// Register the raw paste handler
	router.With(server.v2MiddlewareInjectPaste).Get("/{paste_id}/raw", func(writer http.ResponseWriter, request *http.Request) {
		paste, ok := request.Context().Value("paste").(*pastes.Paste)
		if !ok {
			writeString(writer, http.StatusInternalServerError, "missing paste object")
			return
		}
		_, _ = writer.Write([]byte(paste.Content))
	})

	// Register the paste API endpoints
	router.Get("/api/*", router.NotFoundHandler())
	router.With(server.v2MiddlewareInjectPaste).Get("/api/v2/pastes/{paste_id}", server.v2EndpointGetPaste)
	router.Post("/api/v2/pastes", server.v2EndpointCreatePaste)
	router.With(server.v2MiddlewareInjectPaste, server.v2MiddlewareAuthorize).Patch("/api/v2/pastes/{paste_id}", server.v2EndpointModifyPaste)
	router.With(server.v2MiddlewareInjectPaste, server.v2MiddlewareAuthorize).Delete("/api/v2/pastes/{paste_id}", server.v2EndpointDeletePaste)
	if server.ReportClient != nil {
		router.With(server.v2MiddlewareInjectPaste).Post("/api/v2/pastes/{paste_id}/report", server.v2EndpointReportPaste)
	}
	router.Get("/api/v2/info", func(writer http.ResponseWriter, request *http.Request) {
		writeJSONOrErr(request, writer, http.StatusOK, map[string]any{
			"version":            meta.Version,
			"modificationTokens": server.ModificationTokensEnabled,
			"reports":            server.ReportClient != nil,
			"pasteLifetime":      -1, // TODO: Return paste lifetime
		})
	})

	// Start the HTTP server
	server.httpServer = &http.Server{
		Addr:    server.Address,
		Handler: router,
	}
	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	if err := server.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	server.httpServer = nil
	return nil
}
