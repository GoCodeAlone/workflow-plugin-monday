package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://api.monday.com/v2"
	defaultAPIVersion = "2026-04"
)

type MondayClient struct {
	baseURL    string
	apiToken   string
	apiVersion string
	httpClient *http.Client
}

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLResponse struct {
	Data       json.RawMessage `json:"data"`
	Errors     []graphQLError  `json:"errors,omitempty"`
	Complexity *complexity     `json:"complexity,omitempty"`
}

type graphQLError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
}

type complexity struct {
	Before int `json:"before"`
	After  int `json:"after"`
	Query  int `json:"query"`
}

type ClientOption func(*MondayClient)

func WithBaseURL(url string) ClientOption {
	return func(c *MondayClient) { c.baseURL = url }
}

func WithAPIVersion(version string) ClientOption {
	return func(c *MondayClient) { c.apiVersion = version }
}

func NewMondayClient(apiToken string, opts ...ClientOption) *MondayClient {
	c := &MondayClient{
		baseURL:    defaultBaseURL,
		apiToken:   apiToken,
		apiVersion: defaultAPIVersion,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *MondayClient) Execute(ctx context.Context, query string, variables map[string]any) (json.RawMessage, error) {
	body, err := json.Marshal(graphQLRequest{Query: query, Variables: variables})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiToken)
	req.Header.Set("API-Version", c.apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited (HTTP 429)")
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var gqlResp graphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	return gqlResp.Data, nil
}

func (c *MondayClient) ExecuteInto(ctx context.Context, query string, variables map[string]any, target any) error {
	data, err := c.Execute(ctx, query, variables)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
