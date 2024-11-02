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

type AstroClient interface {
	GetByCity(ctx context.Context, city string) (*dto.AstroByCity, error)
}

type AstroAPIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewAstroAPIClient(openWeatherBaseURL, openWeatherApiKey string) AstroClient {
	cl := NewHttpClient(WithTimeout(3 * time.Second))
	return &AstroAPIClient{
		baseURL: openWeatherBaseURL,
		apiKey:  openWeatherApiKey,
		client:  cl,
	}
}

func (api *AstroAPIClient) GetByCity(ctx context.Context, city string) (*dto.AstroByCity, error) {

	astro, err := api.httpGetByCity(ctx, city)

	if err != nil {
		var apiErr AstroApiErr
		ok := errors.As(err, &apiErr)
		if ok && apiErr.Cod == 404 {
			return nil, result.NotFoundErr(apiErr.Error())
		}
		return nil, result.InternalServerErr("Failed to fetch astronomy data: " + err.Error())
	}

	return astro, nil
}

func (api *AstroAPIClient) httpGetByCity(ctx context.Context, city string) (*dto.AstroByCity, error) {

	encodedCity := url.QueryEscape(city)
	endpoint := fmt.Sprintf("/data/2.5/weather?q=%s&appid=%s", encodedCity, api.apiKey)

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
		msg := "failed to get astro by city"
		slog.Error(msg, slog.String("city", city), slog.String("status", response.Status))

		var apiErr AstroApiErr
		err = json.NewDecoder(response.Body).Decode(&apiErr)
		if err != nil {
			return nil, apiErr
		}

		return nil, errors.New(msg)
	}

	var astro dto.AstroByCity
	err = json.NewDecoder(response.Body).Decode(&astro)
	if err != nil {
		return nil, err
	}

	return &astro, nil
}

// AstroApiErr represents an error response from the OpenWeather API.
type AstroApiErr struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

func (w AstroApiErr) Error() string {
	return w.Message
}
