package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

type APIClientOption func(*APIClient)

func NewAPIClient(baseURL string, o ...APIClientOption) *APIClient {
	defaultCl := &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range o {
		opt(defaultCl)
	}

	return defaultCl
}

func WithTimeout(timeout time.Duration) APIClientOption {
	return func(c *APIClient) {
		c.httpClient.Timeout = timeout
	}
}

func (c *APIClient) Get(ctx context.Context, endpoint string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		msg := "failed to get data from external API"
		slog.Error(msg, slog.String("status", resp.Status))
		return errors.New(msg)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
