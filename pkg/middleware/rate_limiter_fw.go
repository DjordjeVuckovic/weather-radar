package middleware

import (
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type FixedWindowLimiterConfig struct {
	Window      time.Duration
	MaxRequests int32
}

type FixedWindowLimiter struct {
	window      time.Duration
	maxRequests int32
	clients     sync.Map
}

func NewFixedWindowLimiter(config FixedWindowLimiterConfig) *FixedWindowLimiter {
	limiter := &FixedWindowLimiter{
		window:      config.Window,
		maxRequests: config.MaxRequests,
	}
	return limiter
}

func (fw *FixedWindowLimiter) AddAndCheckLimit(r *http.Request) (Limit, error) {
	clientID := getClientID(r)

	cl, _ := fw.clients.LoadOrStore(clientID, &clientLimit{
		requestCount: 0,
	})
	client := cl.(*clientLimit)

	now := time.Now()
	currentWindowStart := client.windowStart.Load()

	if now.UnixNano()-currentWindowStart >= fw.window.Nanoseconds() {
		client.windowStart.Store(now.UnixNano())
		atomic.StoreInt32(&client.requestCount, 1)
	} else {
		atomic.AddInt32(&client.requestCount, 1)
	}

	requestCount := atomic.LoadInt32(&client.requestCount)
	limitExceeded := requestCount > fw.maxRequests
	remaining := fw.maxRequests - requestCount

	if remaining < 0 {
		remaining = 0
	}
	resetTime := time.Unix(0, client.windowStart.Load()).Add(fw.window)
	return Limit{
		Exceeded:  limitExceeded,
		Limit:     int(fw.maxRequests),
		Remaining: int(remaining),
		Reset:     resetTime,
	}, nil
}
