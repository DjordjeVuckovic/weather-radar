package client

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"time"
)

type MockWeatherClient struct {
	Response *dto.WeatherByCity
	Error    error
	Delay    time.Duration
}

func NewMockWeatherClient(err error, delay time.Duration) *MockWeatherClient {
	return &MockWeatherClient{
		Response: &dto.WeatherByCity{
			Location: struct {
				Name           string  `json:"name"`
				Region         string  `json:"region"`
				Country        string  `json:"country"`
				Lat            float64 `json:"lat"`
				Lon            float64 `json:"lon"`
				TzId           string  `json:"tz_id"`
				LocaltimeEpoch int     `json:"localtime_epoch"`
				Localtime      string  `json:"localtime"`
			}{
				Name:           "London",
				Region:         "City of London, Greater London",
				Country:        "United Kingdom",
				Lat:            12,
				Lon:            12,
				TzId:           "Europe/London",
				LocaltimeEpoch: 1730668035,
				Localtime:      "12:00",
			},
		},
		Error: err,
		Delay: delay,
	}
}

func (m *MockWeatherClient) GetByCity(ctx context.Context, _ string) (*dto.WeatherByCity, error) {
	if m.Delay > 0 {
		select {
		case <-time.After(m.Delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return m.Response, m.Error
}

type MockAstroClient struct {
	Response *dto.AstroByCity
	Error    error
}

func NewMockAstroClient(err error) *MockAstroClient {
	return &MockAstroClient{
		Response: &dto.AstroByCity{
			Sys: struct {
				Type    int    `json:"type"`
				Id      int    `json:"id"`
				Country string `json:"country"`
				Sunrise int    `json:"sunrise"`
				Sunset  int    `json:"sunset"`
			}{Type: 1, Id: 1, Country: "SRB", Sunrise: 1730668035, Sunset: 1730668035},
		},
		Error: err,
	}
}

func (m *MockAstroClient) GetByCity(ctx context.Context, city string) (*dto.AstroByCity, error) {
	return m.Response, m.Error
}
