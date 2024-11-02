package client

import (
	"net/http"
	"time"
)

type HttpClientOption func(*http.Client)

func NewHttpClient(o ...HttpClientOption) *http.Client {
	defaultCl := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, opt := range o {
		opt(defaultCl)
	}

	return defaultCl
}

func WithTimeout(timeout time.Duration) HttpClientOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}
