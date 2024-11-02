package usecase

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
)

func WeatherByCityQuery(ctx context.Context, client client.WeatherClient, city string) (*dto.WeatherByCity, error) {
	return client.GetWeatherByCity(ctx, city)
}
