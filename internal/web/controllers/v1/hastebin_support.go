package v1

import (
	"encoding/json"
	"time"

	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/shared"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/internal/utils"
	"github.com/valyala/fasthttp"
)

// HastebinSupportHandler handles the legacy hastebin requests
func HastebinSupportHandler(ctx *fasthttp.RequestCtx) {
	// Check content length before reading body into memory
	if config.Current.LengthCap > 0 &&
		ctx.Request.Header.ContentLength() > config.Current.LengthCap {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("request body length overflow")
		return
	}

	// Define the paste content
	var content string
	switch string(ctx.Request.Header.ContentType()) {
	case "text/plain":
		content = string(ctx.PostBody())
	case "multipart/form-data":
		content = string(ctx.FormValue("data"))
	default:
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("invalid content type")
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
		Content:    content,
		Created:    time.Now().Unix(),
		AutoDelete: config.Current.AutoDelete.Enabled,
	}

	// Set a deletion token
	if config.Current.DeletionTokens {
		paste.DeletionToken = utils.RandomString(config.Current.DeletionTokenLength)

		err = paste.HashDeletionToken()
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

	// Respond with the paste key
	jsonData, _ := json.Marshal(map[string]string{
		"key": paste.ID,
	})
	ctx.SetBody(jsonData)
}
