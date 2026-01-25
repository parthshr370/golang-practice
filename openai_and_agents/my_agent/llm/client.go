package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(apikey string) *Client {
	return &Client{
		APIKey:     apikey,
		BaseURL:    "https://openrouter.ai/api/v1",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) CreateChat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {

	// this is essentially converting the request to json for Marshalling
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal Data here please check again %w", err)
	}

	// request the url with all the elements
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Unable to create the request %w ", err)

	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch response check your API %w ", err)

	}

	defer resp.Body.Close() // close the flowing pipe you just stareted

	if resp.StatusCode != http.StatusOK {
		// this is good practice Read the error body to see why failed (optional but good practice)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

	}

	var chatResp ChatResponse
	// this tells that you can just take the response body and the point it to the chatresponse in memory with obviously ChatResponse struct
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &chatResp, nil
}
