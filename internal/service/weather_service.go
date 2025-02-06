package service

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
	"github.com/DjordjeVuckovic/weather-radar/internal/storage"
	"sync"
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

	weatherCh := make(chan *dto.WeatherByCity, 1)
	astroCh := make(chan *dto.AstroByCity, 1)
	errCh := make(chan error, 1)

	go func() {
		wth, err := w.weatherClient.GetByCity(timeoutCtx, city)
		if err != nil {
			errCh <- err
			return
		}
		select {
		case weatherCh <- wth:
		case <-timeoutCtx.Done():
			return

		}
	}()

	go func() {
		a, err := w.astroClient.GetByCity(timeoutCtx, city)
		if err != nil {
			errCh <- err
			return
		}
		select {
		case astroCh <- a:
		case <-timeoutCtx.Done():
			return

		}
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

type AggregatedWeather struct {
	Weather *model.Weather
	City    string
}

func (w *WeatherService) GetWeatherByCites(ctx context.Context, cities []string) ([]AggregatedWeather, error) {
	var (
		resultCh = make(chan AggregatedWeather, len(cities))
		errCh    = make(chan error, len(cities))
		wg       sync.WaitGroup
	)

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			weather, err := w.GetWeatherByCity(ctx, city)
			if err != nil {
				errCh <- err
				return
			}
			resultCh <- AggregatedWeather{
				Weather: weather,
				City:    city,
			}
		}(city)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	var weathers []AggregatedWeather
	for {
		select {
		case result, ok := <-resultCh:
			if !ok {
				return weathers, nil
			}
			weathers = append(weathers, result)
		case err := <-errCh:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (w *WeatherService) GetWeatherStreamByCities(ctx context.Context, cities []string) (<-chan AggregatedWeather, <-chan error) {
	resultCh := make(chan AggregatedWeather, 1)
	errCh := make(chan error, len(cities))

	go func() {
		defer close(resultCh)
		defer close(errCh)
		var wg sync.WaitGroup

		for _, city := range cities {
			wg.Add(1)
			go func(city string) {
				defer wg.Done()
				weather, err := w.GetWeatherByCity(ctx, city)
				if err != nil {
					select {
					case errCh <- err:
					case <-ctx.Done():
					}
					return
				}
				select {
				case resultCh <- AggregatedWeather{Weather: weather, City: city}:
				case <-ctx.Done():
					return
				}
			}(city)
		}

		wg.Wait()
	}()

	return resultCh, errCh
}

func (w *WeatherService) SubmitFeedback(ctx context.Context, feedback *dto.WeatherFeedbackReq) error {
	fb := model.NewFeedbackFromDto(feedback)
	if err := w.storage.AddFeedback(ctx, fb); err != nil {
		return err
	}
	return nil
}
