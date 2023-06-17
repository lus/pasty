package reports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Report struct {
	Paste  string `json:"paste"`
	Reason string `json:"reason"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Client struct {
	WebhookURL   string
	WebhookToken string
}

func (client *Client) Send(report *Report) (*Response, error) {
	data, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, client.WebhookURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if client.WebhookToken != "" {
		request.Header.Set("Authorization", "Bearer "+client.WebhookToken)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("the report webhook responded with an unexpected error: %d (%s)", response.StatusCode, string(body))
	}

	reportResponse := new(Response)
	if err := json.Unmarshal(body, &reportResponse); err != nil {
		return nil, err
	}
	return reportResponse, nil
}
