package storage

import (
	"context"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
	"github.com/google/uuid"
	"sync"
)

type WeatherMemStorage struct {
	feedbacks map[uuid.UUID]*model.Feedback
	mx        sync.RWMutex
}

func NewWeatherInMemStorage() WeatherStorage {
	return &WeatherMemStorage{
		feedbacks: make(map[uuid.UUID]*model.Feedback),
	}
}

func (wms *WeatherMemStorage) AddFeedback(_ context.Context, fb *model.Feedback) error {
	wms.mx.Lock()
	defer wms.mx.Unlock()
	wms.feedbacks[fb.ID] = fb
	return nil
}
