package v2

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/fasthttp/router"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/shared"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/internal/utils"
	limitFasthttp "github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	"github.com/valyala/fasthttp"
)

// InitializePastesController initializes the '/v2/pastes/*' controller
func InitializePastesController(group *router.Group, rateLimiterMiddleware *limitFasthttp.Middleware) {
	// moms spaghetti
	group.GET("/{id}", rateLimiterMiddleware.Handle(middlewareInjectPaste(endpointGetPaste)))
	group.POST("/", rateLimiterMiddleware.Handle(endpointCreatePaste))
	group.DELETE("/{id}", rateLimiterMiddleware.Handle(middlewareInjectPaste(middlewareValidateModificationToken(endpointDeletePaste))))
}

// middlewareInjectPaste retrieves and injects the paste with the specified ID
func middlewareInjectPaste(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		pasteID := ctx.UserValue("id").(string)

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

		ctx.SetUserValue("_paste", paste)

		next(ctx)
	}
}

// middlewareValidateModificationToken extracts and validates a given modification token for an injected paste
func middlewareValidateModificationToken(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		paste := ctx.UserValue("_paste").(*shared.Paste)

		authHeaderSplit := strings.SplitN(string(ctx.Request.Header.Peek("Authorization")), " ", 2)
		if len(authHeaderSplit) < 2 || authHeaderSplit[0] != "Bearer" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("unauthorized")
			return
		}

		modificationToken := authHeaderSplit[1]
		if config.Current.ModificationTokenMaster != "" && modificationToken == config.Current.ModificationTokenMaster {
			next(ctx)
			return
		}
		valid := paste.CheckModificationToken(modificationToken)
		if !valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString("unauthorized")
			return
		}

		next(ctx)
	}
}

// endpointGetPaste handles the 'GET /v2/pastes/{id}' endpoint
func endpointGetPaste(ctx *fasthttp.RequestCtx) {
	paste := ctx.UserValue("_paste").(*shared.Paste)
	paste.DeletionToken = ""
	paste.ModificationToken = ""

	jsonData, err := json.Marshal(paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetBody(jsonData)
}

type endpointCreatePastePayload struct {
	Content string `json:"content"`
}

// endpointCreatePaste handles the 'POST /v2/pastes' endpoint
func endpointCreatePaste(ctx *fasthttp.RequestCtx) {
	// Read, parse and validate the request payload
	payload := new(endpointCreatePastePayload)
	if err := json.Unmarshal(ctx.PostBody(), payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	if payload.Content == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("missing paste content")
		return
	}
	if config.Current.LengthCap > 0 && len(payload.Content) > config.Current.LengthCap {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("too large paste content")
		return
	}

	// Acquire a new paste ID
	id, err := storage.AcquireID()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Prepare the paste object
	paste := &shared.Paste{
		ID:         id,
		Content:    payload.Content,
		Created:    time.Now().Unix(),
		AutoDelete: config.Current.AutoDelete.Enabled,
	}

	// Create a new modification token if enabled
	modificationToken := ""
	if config.Current.ModificationTokens {
		modificationToken = utils.RandomString(config.Current.ModificationTokenLength)
		paste.ModificationToken = modificationToken

		err = paste.HashModificationToken()
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(err.Error())
			return
		}
	}

	// Save the paste
	err = storage.Current.Save(paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Respond with the paste
	pasteCopy := *paste
	pasteCopy.ModificationToken = modificationToken
	jsonData, err := json.Marshal(pasteCopy)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(jsonData)
}

// endpointDeletePaste handles the 'DELETE /v2/pastes/{id}' endpoint
func endpointDeletePaste(ctx *fasthttp.RequestCtx) {
	paste := ctx.UserValue("_paste").(*shared.Paste)
	if err := storage.Current.Delete(paste.ID); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
}
