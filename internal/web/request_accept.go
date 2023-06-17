package web

import "net/http"

func accept(writer http.ResponseWriter, request *http.Request, contentTypes ...string) bool {
	contentType := request.Header.Get("Content-Type")
	for _, accepted := range contentTypes {
		if contentType == accepted {
			return true
		}
	}
	writeString(writer, http.StatusUnsupportedMediaType, "unsupported media type")
	return false
}
