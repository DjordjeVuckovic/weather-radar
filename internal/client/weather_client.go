package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type WeatherClient interface {
	GetByCity(ctx context.Context, city string) (*dto.WeatherByCity, error)
}

type APIWeatherClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewWeatherAPIClient(weatherBaseURL, weatherApiKey string) WeatherClient {
	cl := NewHttpClient(WithTimeout(3 * time.Second))
	return &APIWeatherClient{
		baseURL: weatherBaseURL,
		apiKey:  weatherApiKey,
		client:  cl,
	}
}

func (api *APIWeatherClient) GetByCity(ctx context.Context, city string) (*dto.WeatherByCity, error) {

	weather, err := api.httpGetByCity(ctx, city)

	if err != nil {
		var apiErr WeatherApiErr
		ok := errors.As(err, &apiErr)
		if ok && apiErr.Err.Code == 1006 {
			return nil, result.NotFoundErr(apiErr.Error())
		}
		return nil, result.InternalServerErr("Failed to fetch weather data: " + err.Error())
	}

	return weather, nil
}

func (api *APIWeatherClient) httpGetByCity(ctx context.Context, city string) (*dto.WeatherByCity, error) {

	encodedCity := url.QueryEscape(city)
	endpoint := fmt.Sprintf("/current.json?key=%s&q=%s&aqi=no", api.apiKey, encodedCity)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api.baseURL+endpoint, nil)

	if err != nil {
		return nil, err
	}

	response, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		msg := "failed to get weather by city"
		slog.Error(msg, slog.String("city", city), slog.String("status", response.Status))
		var apiErr WeatherApiErr
		if err := json.NewDecoder(response.Body).Decode(&apiErr); err != nil {
			return nil, err
		}
		return nil, apiErr
	}

	var weather dto.WeatherByCity
	if err := json.NewDecoder(response.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

// WeatherApiErr
// code: 1006 No matching location found
type WeatherApiErr struct {
	Err struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (w WeatherApiErr) Error() string {
	return w.Err.Message
}
