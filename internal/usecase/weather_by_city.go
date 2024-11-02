package usecase

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
)

func WeatherByCityQuery(
	ctx context.Context,
	wCl client.WeatherClient,
	astroCl client.AstroClient,
	city string) (*model.Weather, error) {

	weatherCh := make(chan *dto.WeatherByCity)
	astroCh := make(chan *dto.AstroByCity)
	errCh := make(chan error, 2)

	go func() {
		w, err := wCl.GetByCity(ctx, city)

		if err != nil {
			errCh <- err
			return
		}

		weatherCh <- w
	}()

	go func() {
		a, err := astroCl.GetByCity(ctx, city)

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
		}
	}

	weather := model.NewWeatherFromDto(weatherData, astroData)
	return weather, nil
}
