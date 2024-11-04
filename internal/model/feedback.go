package model

import (
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/google/uuid"
)

type Feedback struct {
	ID      uuid.UUID
	Date    string
	City    string
	Message string
}

func NewFeedbackFromDto(dto *dto.WeatherFeedbackReq) *Feedback {
	return &Feedback{
		ID:      uuid.New(),
		Date:    dto.Date,
		City:    dto.City,
		Message: dto.Message,
	}
}
