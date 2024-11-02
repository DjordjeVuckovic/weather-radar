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
	Environment   string
	CorsOrigins   string
	WeatherUrl    string
	WeatherApiKey string
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

	wUrl := os.Getenv("WEATHER_URL")
	if wUrl == "" {
		panic("WEATHER_URL is required")
	}

	wApiKey := os.Getenv("WEATHER_API_KEY")
	if wApiKey == "" {
		panic("WEATHER_API_KEY is required")
	}

	return Env{
		Environment:   os.Getenv("ENVIRONMENT"),
		CorsOrigins:   strings.Join(origins, ","),
		WeatherUrl:    wUrl,
		WeatherApiKey: wApiKey,
	}
}