package web

import (
	"encoding/json"
	"path/filepath"
	"strings"

	routing "github.com/fasthttp/router"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/static"
	"github.com/lus/pasty/internal/storage"
	v1 "github.com/lus/pasty/internal/web/controllers/v1"
	v2 "github.com/lus/pasty/internal/web/controllers/v2"
	"github.com/ulule/limiter/v3"
	limitFasthttp "github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/valyala/fasthttp"
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
	raw := rawHandler()
	router.GET("/{path:*}", func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		if !strings.HasPrefix(path, "/api") && (strings.Count(path, "/") == 1 || strings.HasPrefix(path, "/assets")) {
			if strings.HasPrefix(path, "/assets/js/") {
				ctx.SetContentType("text/javascript")
			}
			frontend(ctx)
			return
		} else if strings.HasSuffix(strings.TrimSuffix(path, "/"), "/raw") {
			raw(ctx)
			return
		}
		router.NotFound(ctx)
	})

	// Set up the rate limiter
	rate, err := limiter.NewRateFromFormatted(config.Current.RateLimit)
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
				jsonData, _ := json.Marshal(map[string]interface{}{
					"version":        static.Version,
					"deletionTokens": config.Current.ModificationTokens,
				})
				ctx.SetBody(jsonData)
			})
			v1.InitializePastesController(v1Route.Group("/pastes"), rateLimiterMiddleware)
		}

		v2Route := apiRoute.Group("/v2")
		{
			v2Route.GET("/info", func(ctx *fasthttp.RequestCtx) {
				jsonData, _ := json.Marshal(map[string]interface{}{
					"version":            static.Version,
					"modificationTokens": config.Current.ModificationTokens,
					"reports":            config.Current.Reports.Reports,
				})
				ctx.SetBody(jsonData)
			})
			v2.InitializePastesController(v2Route.Group("/pastes"), rateLimiterMiddleware)
		}
	}

	// Route the hastebin documents route if hastebin support is enabled
	if config.Current.HastebinSupport {
		router.POST("/documents", rateLimiterMiddleware.Handle(v1.HastebinSupportHandler))
	}

	// Serve the web resources
	return (&fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			// Add the CORS headers
			ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")

			// Call the router handler
			router.Handler(ctx)
		},
		Logger: new(nilLogger),
	}).ListenAndServe(config.Current.WebAddress)
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

func rawHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		pathSanitized := strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
		pasteID := strings.TrimSuffix(pathSanitized, "/raw")

		paste, err := storage.Current.Get(pasteID)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(err.Error())
			return
		}

		if paste == nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBodyString("paste not found")
			return
		}

		ctx.SetBodyString(paste.Content)
	}
}
