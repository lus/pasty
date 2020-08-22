package v1

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"github.com/bwmarrin/snowflake"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// InitializePastesController initializes the '/v1/pastes/*' controller
func InitializePastesController(group *router.Group) {
	group.GET("/{id}", v1GetPaste)
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

	// Respond with the paste
	jsonData, err := json.Marshal(paste)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
		return
	}
	ctx.SetBody(jsonData)
}
