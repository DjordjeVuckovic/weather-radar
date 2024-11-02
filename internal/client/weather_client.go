package client

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"time"
)

type WeatherClient interface {
	GetWeatherByCity(ctx context.Context, city string) (*dto.WeatherByCity, error)
}

type APIWeatherClient struct {
	baseURL string
	apiKey  string
}

func NewAPIWeatherClient(baseURL, apiKey string) WeatherClient {
	return &APIWeatherClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (api *APIWeatherClient) GetWeatherByCity(ctx context.Context, city string) (*dto.WeatherByCity, error) {
	cl := NewAPIClient(
		api.baseURL,
		WithTimeout(5*time.Second),
	)

	var weather *dto.WeatherByCity
	endpoint := "/current.json" + "?key=" + api.apiKey + "&q=" + city + "&aqi=no"
	err := cl.Get(ctx, endpoint, weather)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
