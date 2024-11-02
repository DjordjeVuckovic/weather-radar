package config

import (
	"log/slog"
	"os"
	"strings"
)
import (
	"github.com/joho/godotenv"
)

type Env struct {
	ENV               string
	CorsOrigins       string
	WeatherUrl        string
	WeatherApiKey     string
	OpenWeatherUrl    string
	OpenWeatherApiKey string
}

func Load() Env {
	if err := godotenv.Load(); err != nil {
		slog.Info("Skipping .config file ...")
	}
	origin := os.Getenv("CORS_ORIGINS")
	origins := strings.Split(origin, ",")
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	wUrl := os.Getenv("WEATHER_API_URL")
	if wUrl == "" {
		panic("WEATHER_API_URL is required")
	}

	wApiKey := os.Getenv("WEATHER_API_KEY")
	if wApiKey == "" {
		panic("WEATHER_API_KEY is required")
	}

	owUrl := os.Getenv("OPEN_WEATHER_API_URL")
	if owUrl == "" {
		panic("OPEN_WEATHER_API_URL is required")
	}

	owApiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if owApiKey == "" {
		panic("OPEN_WEATHER_API_KEY is required")
	}

	return Env{
		ENV:               os.Getenv("ENV"),
		CorsOrigins:       strings.Join(origins, ","),
		WeatherUrl:        wUrl,
		WeatherApiKey:     wApiKey,
		OpenWeatherUrl:    owUrl,
		OpenWeatherApiKey: owApiKey,
	}
}
