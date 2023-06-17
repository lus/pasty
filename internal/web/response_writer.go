package web

import (
	"encoding/json"
	"github.com/lus/pasty/pkg/chizerolog"
	"net/http"
	"strconv"
)

func writeErr(request *http.Request, writer http.ResponseWriter, err error) {
	chizerolog.InjectError(request, err)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Set("Content-Length", strconv.Itoa(len(err.Error())))
	writeString(writer, http.StatusInternalServerError, err.Error())
}

func writeString(writer http.ResponseWriter, status int, value string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Set("Content-Length", strconv.Itoa(len(value)))
	writer.WriteHeader(status)
	writer.Write([]byte(value))
}

func writeJSON(writer http.ResponseWriter, status int, value any) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(jsonData)))
	writer.WriteHeader(status)
	writer.Write(jsonData)

	return nil
}

func writeJSONOrErr(request *http.Request, writer http.ResponseWriter, status int, value any) {
	if err := writeJSON(writer, status, value); err != nil {
		writeErr(request, writer, err)
	}
}
