package web

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/env"
	"github.com/Lukaesebrot/pasty/internal/static"
	v1 "github.com/Lukaesebrot/pasty/internal/web/controllers/v1"
	routing "github.com/fasthttp/router"
	"github.com/ulule/limiter/v3"
	limitFasthttp "github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/valyala/fasthttp"
	"path/filepath"
	"strings"
)

// Serve serves the web resources
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
		if !strings.HasPrefix(path, "/api") && (strings.Count(path, "/") == 1 || strings.HasPrefix(path, "/assets")) {
			frontend(ctx)
			return
		}
		router.NotFound(ctx)
	})

	// Set up the rate limiter
	rate, err := limiter.NewRateFromFormatted(env.Get("RATE_LIMIT", "30-M"))
	if err != nil {
		return err
	}
	rateLimiter := limiter.New(memory.NewStore(), rate)
	rateLimiterMiddleware := limitFasthttp.NewMiddleware(rateLimiter)

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
			v1.InitializePastesController(v1Route.Group("/pastes"), rateLimiterMiddleware)
		}
	}

	// Route the hastebin documents route if hastebin support is enabled
	if env.Bool("HASTEBIN_SUPPORT", false) {
		router.POST("/documents", rateLimiterMiddleware.Handle(v1.HastebinSupportHandler))
	}

	// Serve the web resources
	address := env.Get("WEB_ADDRESS", ":8080")
	return (&fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			// Add the CORS headers
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

			// Call the router handler
			router.Handler(ctx)
		},
		Logger: new(nilLogger),
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
		if strings.HasPrefix(string(ctx.Path()), "/assets") {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("not found")
			return
		}
		ctx.SendFile(filepath.Join(fs.Root, "index.html"))
	}
	return fs.NewRequestHandler()
}
