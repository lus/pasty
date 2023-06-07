package web

import (
	"encoding/json"
	"net/http"
)

func writeErr(writer http.ResponseWriter, err error) {
	writeString(writer, http.StatusInternalServerError, err.Error())
}

func writeString(writer http.ResponseWriter, status int, value string) {
	writer.WriteHeader(status)
	writer.Write([]byte(value))
}

func writeJSON(writer http.ResponseWriter, status int, value any) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	writer.WriteHeader(status)
	writer.Write(jsonData)

	return nil
}

func writeJSONOrErr(writer http.ResponseWriter, status int, value any) {
	if err := writeJSON(writer, status, value); err != nil {
		writeErr(writer, err)
	}
}
