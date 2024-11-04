package storage

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
)

type WeatherStorage interface {
	AddFeedback(ctx context.Context, fb *model.Feedback) error
}
