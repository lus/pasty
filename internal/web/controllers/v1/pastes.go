package v1

import (
	"encoding/json"
	"time"

	"github.com/fasthttp/router"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/shared"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/internal/utils"
	limitFasthttp "github.com/ulule/limiter/v3/drivers/middleware/fasthttp"
	"github.com/valyala/fasthttp"
)

// InitializePastesController initializes the '/v1/pastes/*' controller
func InitializePastesController(group *router.Group, rateLimiterMiddleware *limitFasthttp.Middleware) {
	group.GET("/{id}", rateLimiterMiddleware.Handle(v1GetPaste))
	group.POST("", rateLimiterMiddleware.Handle(v1PostPaste))
	group.DELETE("/{id}", rateLimiterMiddleware.Handle(v1DeletePaste))
}

// v1GetPaste handles the 'GET /v1/pastes/{id}' endpoint
func v1GetPaste(ctx *fasthttp.RequestCtx) {
	// Read the ID
	id := ctx.UserValue("id").(string)

	// Retrieve the paste
	paste, err := storage.Current.Get(id)
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
	legacyPaste := legacyFromModern(paste)
	legacyPaste.DeletionToken = ""

	// Respond with the paste
	jsonData, err := json.Marshal(legacyPaste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetBody(jsonData)
}

// v1PostPaste handles the 'POST /v1/pastes' endpoint
func v1PostPaste(ctx *fasthttp.RequestCtx) {
	// Check content length before reading body into memory
	if config.Current.LengthCap > 0 &&
		ctx.Request.Header.ContentLength() > config.Current.LengthCap {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("request body length overflow")
		return
	}

	// Unmarshal the body
	values := make(map[string]string)
	err := json.Unmarshal(ctx.PostBody(), &values)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid request body")
		return
	}

	// Validate the content of the paste
	if values["content"] == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("missing 'content' field")
		return
	}

	// Acquire the paste ID
	id, err := storage.AcquireID()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Create the paste object
	paste := &shared.Paste{
		ID:         id,
		Content:    values["content"],
		Created:    time.Now().Unix(),
		AutoDelete: config.Current.AutoDelete.Enabled,
	}

	// Set a modification token
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
	pasteCopy := legacyFromModern(paste)
	pasteCopy.DeletionToken = modificationToken
	jsonData, err := json.Marshal(pasteCopy)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetBody(jsonData)
}

// v1DeletePaste handles the 'DELETE /v1/pastes/{id}'
func v1DeletePaste(ctx *fasthttp.RequestCtx) {
	// Read the ID
	id := ctx.UserValue("id").(string)

	// Unmarshal the body
	values := make(map[string]string)
	err := json.Unmarshal(ctx.PostBody(), &values)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid request body")
		return
	}

	// Validate the modification token of the paste
	modificationToken := values["deletionToken"]
	if modificationToken == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("missing 'deletionToken' field")
		return
	}

	// Retrieve the paste
	paste, err := storage.Current.Get(id)
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

	// Check if the modification token is correct
	if (config.Current.ModificationTokenMaster == "" || modificationToken != config.Current.ModificationTokenMaster) && !paste.CheckModificationToken(modificationToken) {
		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.SetBodyString("invalid deletion token")
		return
	}

	// Delete the paste
	err = storage.Current.Delete(paste.ID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Respond with 'ok'
	ctx.SetBodyString("ok")
}
