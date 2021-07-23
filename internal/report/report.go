package report

import (
	"encoding/json"
	"fmt"

	"github.com/lus/pasty/internal/config"
	"github.com/valyala/fasthttp"
)

// ReportRequest represents a report request sent to the report webhook
type ReportRequest struct {
	Paste     string `json:"paste"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

// ReportResponse represents a report response received from the report webhook
type ReportResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SendReport sends a report request to the report webhook
func SendReport(reportRequest *ReportRequest) (*ReportResponse, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	request.Header.SetMethod(fasthttp.MethodPost)
	request.SetRequestURI(config.Current.Reports.ReportWebhook)
	if config.Current.Reports.ReportWebhookToken != "" {
		request.Header.Set("Authorization", "Bearer "+config.Current.Reports.ReportWebhookToken)
	}

	data, err := json.Marshal(reportRequest)
	if err != nil {
		return nil, err
	}
	request.SetBody(data)

	if err := fasthttp.Do(request, response); err != nil {
		return nil, err
	}

	status := response.StatusCode()
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("the report webhook responded with an unexpected error: %d (%s)", status, string(response.Body()))
	}

	reportResponse := new(ReportResponse)
	if err := json.Unmarshal(response.Body(), reportResponse); err != nil {
		return nil, err
	}
	return reportResponse, nil
}
