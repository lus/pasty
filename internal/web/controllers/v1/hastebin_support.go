package v1

import (
	"encoding/json"
	"github.com/Lukaesebrot/pasty/internal/pastes"
	"github.com/Lukaesebrot/pasty/internal/storage"
	"github.com/valyala/fasthttp"
)

// HastebinSupportHandler handles the legacy hastebin requests
func HastebinSupportHandler(ctx *fasthttp.RequestCtx) {
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
