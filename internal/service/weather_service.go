package service

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
	"github.com/DjordjeVuckovic/weather-radar/internal/storage"
	"time"
)

const timeout = 1 * time.Second

type WeatherService struct {
	weatherClient client.WeatherClient
	astroClient   client.AstroClient
	storage       storage.WeatherStorage
}

func NewWeatherService(wCl client.WeatherClient, aCl client.AstroClient, st storage.WeatherStorage) *WeatherService {
	return &WeatherService{
		weatherClient: wCl,
		astroClient:   aCl,
		storage:       st,
	}
}

func (w *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*model.Weather, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	weatherCh := make(chan *dto.WeatherByCity)
	astroCh := make(chan *dto.AstroByCity)
	errCh := make(chan error, 1)

	go func() {
		wth, err := w.weatherClient.GetByCity(timeoutCtx, city)
		if err != nil {
			errCh <- err
			return
		}

		weatherCh <- wth
	}()

	go func() {
		a, err := w.astroClient.GetByCity(timeoutCtx, city)
		if err != nil {
			errCh <- err
			return
		}

		astroCh <- a
	}()

	var weatherData *dto.WeatherByCity
	var astroData *dto.AstroByCity
	for i := 0; i < 2; i++ {
		select {
		case w := <-weatherCh:
			weatherData = w
		case a := <-astroCh:
			astroData = a
		case err := <-errCh:
			return nil, err
		case <-timeoutCtx.Done():
			return nil, context.DeadlineExceeded
		}
	}

	weather := model.NewWeatherFromDto(weatherData, astroData)
	return weather, nil
}

func (w *WeatherService) SubmitFeedback(ctx context.Context, feedback *dto.WeatherFeedbackReq) error {
	fb := model.NewFeedbackFromDto(feedback)
	if err := w.storage.AddFeedback(ctx, fb); err != nil {
		return err
	}
	return nil
}
