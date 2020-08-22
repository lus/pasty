package web

import (
	"github.com/Lukaesebrot/pasty/internal/env"
	v1 "github.com/Lukaesebrot/pasty/internal/web/controllers/v1"
	routing "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Serve serves the web server
func Serve() error {
	// Create the router
	router := routing.New()

	// Route the API endpoints
	apiRoute := router.Group("/api")
	{
		v1Route := apiRoute.Group("/v1")
		{
			v1.InitializePastesController(v1Route.Group("/pastes"))
		}
	}

	// TODO: Route the paste endpoints

	// Serve the web server
	address := env.Get("WEB_ADDRESS", ":8080")
	return fasthttp.ListenAndServe(address, router.Handler)
}
