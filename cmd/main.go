package main

import (
	"github.com/DjordjeVuckovic/weather-radar/api"
	"github.com/DjordjeVuckovic/weather-radar/internal/cache"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/config"
	"github.com/DjordjeVuckovic/weather-radar/internal/service"
	"github.com/DjordjeVuckovic/weather-radar/pkg/logger"
	"github.com/DjordjeVuckovic/weather-radar/pkg/middleware"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"time"
)

func main() {
	cfg := config.Load()

	logger.InitSlog(logger.Config{
		Level:   logger.DebugLevel,
		Handler: logger.Text,
	})

	gst := server.WithGracefulShutdownTimeout(5 * time.Second)
	s := server.NewServer(":1312", gst)

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(middleware.CORS(middleware.Config{Origin: cfg.CorsOrigins}))

	c := cache.NewInMemCache(5*time.Minute, cache.EvictLRU)

	wCl := client.NewWeatherAPIClient(
		cfg.WeatherUrl,
		cfg.WeatherApiKey,
	)
	astroCl := client.NewAstroAPIClient(
		cfg.OpenWeatherUrl,
		cfg.OpenWeatherApiKey,
	)
	authService := service.NewAuthService(service.AuthCredentials{
		Username: cfg.BasicAuthUsername,
		Password: cfg.BasicAuthPassword,
	})

	wService := service.NewWeatherService(wCl, astroCl)
	api.BindWeatherApi(s, wService, authService, c)

	s.SetupNotFoundHandler()

	if err := s.Start(); err != nil {
		slog.Error(err.Error())
	}
}
