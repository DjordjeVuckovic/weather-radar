package service

import (
	"context"
	"errors"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/internal/storage"
	"testing"
	"time"
)

func TestGetWeatherByCity_Success(t *testing.T) {
	weatherMock := client.NewMockWeatherClient(nil, 0)
	astroMock := client.NewMockAstroClient(nil)

	service := NewWeatherService(weatherMock, astroMock, storage.NewWeatherInMemStorage())
	ctx := context.Background()

	city := "London"
	weather, err := service.GetWeatherByCity(ctx, city)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if weather == nil {
		t.Fatal("Expected non-nil weather response")
	}
	if weather.Location.Name != city {
		t.Fatalf("Expected weather data with City: London  got %+v", weather)
	}
}

func TestGetWeatherByCity_WeatherClientError(t *testing.T) {
	errMsg := "no matching location found"
	weatherMock := client.NewMockWeatherClient(errors.New(errMsg), 0)
	astroMock := client.NewMockAstroClient(errors.New(errMsg))

	service := NewWeatherService(weatherMock, astroMock, storage.NewWeatherInMemStorage())
	ctx := context.Background()

	weather, err := service.GetWeatherByCity(ctx, "London")

	if err == nil || err.Error() != errMsg {
		t.Fatalf("Expected weather error, got %v", err)
	}
	if weather != nil {
		t.Fatal("Expected nil weather response due to error")
	}
}

func TestGetWeatherByCity_Timeout(t *testing.T) {
	weatherMock := client.NewMockWeatherClient(nil, 10*time.Millisecond)
	astroMock := client.NewMockAstroClient(nil)

	service := NewWeatherService(weatherMock, astroMock, storage.NewWeatherInMemStorage())
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	weather, err := service.GetWeatherByCity(ctx, "Paris")

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Expected DeadlineExceeded error, got %v", err)
	}
	if weather != nil {
		t.Fatal("Expected nil weather response due to timeout")
	}
}

func TestWeatherService_SubmitFeedback(t *testing.T) {
	st := storage.NewWeatherInMemStorage()
	service := NewWeatherService(nil, nil, st)

	feedback := &dto.WeatherFeedbackReq{
		City:    "London",
		Date:    "2024-10-27",
		Message: "It was raining",
	}
	err := service.SubmitFeedback(context.Background(), feedback)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
