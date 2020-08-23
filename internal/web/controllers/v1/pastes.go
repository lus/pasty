package v1

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"github.com/bwmarrin/snowflake"
	"github.com/fasthttp/router"
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
	// Parse the ID
	id, err := snowflake.ParseString(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid ID format")
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

	// Respond with the paste
	jsonData, err := json.Marshal(paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetBody(jsonData)
}

// v1PostPaste handles the 'POST /v1/pastes' endpoint
func v1PostPaste(ctx *fasthttp.RequestCtx) {
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

	// Create the paste object
	paste, err := pastes.Create(values["content"])
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}

	// Hash the deletion token
	pasteCopy := *paste
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

	// Respond with the paste
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
	// Parse the ID
	id, err := snowflake.ParseString(ctx.UserValue("id").(string))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid ID format")
		return
	}

	// Unmarshal the body
	values := make(map[string]string)
	err = json.Unmarshal(ctx.PostBody(), &values)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid request body")
		return
	}

	// Validate the deletion token of the paste
	if values["deletionToken"] == "" {
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

	// Check if the deletion token is correct
	if !paste.CheckDeletionToken(values["deletionToken"]) {
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
