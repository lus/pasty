package web

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/Lukaesebrot/pasty/internal/static"
	"github.com/Lukaesebrot/pasty/internal/storage"
	v1 "github.com/Lukaesebrot/pasty/internal/web/controllers/v1"
	routing "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strings"
)

// Serve serves the web server
func Serve() error {
	// Create the router
	router := routing.New()

	// Define the 404 handler
	router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString("not found")
	}

	// Route the frontend requests
	frontend := frontendHandler()
	router.GET("/{path:*}", func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if !strings.HasPrefix(path, "/api") {
			frontend(ctx)
			return
		}
		router.NotFound(ctx)
	})

	// Route the API endpoints
	apiRoute := router.Group("/api")
	{
		v1Route := apiRoute.Group("/v1")
		{
			v1Route.GET("/info", func(ctx *fasthttp.RequestCtx) {
				jsonData, _ := json.Marshal(map[string]string{
					"version": static.Version,
				})
				ctx.SetBody(jsonData)
			})
			v1.InitializePastesController(v1Route.Group("/pastes"))
		}
	}

	// Route the hastebin documents route if hastebin support is enabled
	if env.Get("HASTEBIN_SUPPORT", "false") == "true" {
		router.POST("/documents", hastebinSupportHandler)
	}

	// Serve the web server
	address := env.Get("WEB_ADDRESS", ":8080")
	return (&fasthttp.Server{
		Handler: router.Handler,
		Logger:  new(nilLogger),
	}).ListenAndServe(address)
}

// frontendHandler handles the frontend routing
func frontendHandler() fasthttp.RequestHandler {
	// Create the file server
	fs := &fasthttp.FS{
		Root:          "./web",
		IndexNames:    []string{"index.html"},
		CacheDuration: 0,
	}
	fs.PathNotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SendFile(filepath.Join(fs.Root, "index.html"))
	}
	return fs.NewRequestHandler()
}

// hastebinSupportHandler handles the legacy hastebin requests
func hastebinSupportHandler(ctx *fasthttp.RequestCtx) {
	// Define the paste content
	var content string
	switch string(ctx.Request.Header.ContentType()) {
	case "text/plain":
		content = string(ctx.PostBody())
		break
	case "multipart/form-data":
		content = string(ctx.FormValue("data"))
		break
	default:
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid content type")
		return
	}

	// Create the paste object
	paste, err := pastes.Create(content)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Hash the deletion token
	err = paste.HashDeletionToken()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Save the paste
	err = storage.Current.Save(paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Respond with the paste key
	jsonData, _ := json.Marshal(map[string]string{
		"key": paste.ID.String(),
	})
	ctx.SetBody(jsonData)
}
