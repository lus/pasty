package web

import (
	"encoding/json"
	"github.com/lus/pasty/internal/pastes"
	"github.com/lus/pasty/internal/reports"
	"io"
	"net/http"
)

type v2EndpointReportPastePayload struct {
	Reason string `json:"reason"`
}

func (server *Server) v2EndpointReportPaste(writer http.ResponseWriter, request *http.Request) {
	paste, ok := request.Context().Value("paste").(*pastes.Paste)
	if !ok {
		writeString(writer, http.StatusInternalServerError, "missing paste object")
		return
	}

	// Read, parse and validate the request payload
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writeErr(request, writer, err)
		return
	}
	payload := new(v2EndpointReportPastePayload)
	if err := json.Unmarshal(body, payload); err != nil {
		writeErr(request, writer, err)
		return
	}
	if payload.Reason == "" {
		writeString(writer, http.StatusBadRequest, "missing report reason")
		return
	}

	report := &reports.Report{
		Paste:  paste.ID,
		Reason: payload.Reason,
	}
	response, err := server.ReportClient.Send(report)
	if err != nil {
		writeErr(request, writer, err)
		return
	}
	writeJSONOrErr(request, writer, http.StatusOK, response)
}
