package web

import (
	"github.com/Lukaesebrot/pasty/internal/env"
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
			v1.InitializePastesController(v1Route.Group("/pastes"))
		}
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
