package main

import (
	"github.com/DjordjeVuckovic/weather-radar/api"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/config"
	"github.com/DjordjeVuckovic/weather-radar/internal/service"
	"github.com/DjordjeVuckovic/weather-radar/internal/storage"
	"github.com/DjordjeVuckovic/weather-radar/pkg/cache"
	"github.com/DjordjeVuckovic/weather-radar/pkg/logger"
	"github.com/DjordjeVuckovic/weather-radar/pkg/middleware"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"time"
)

func main() {
	cfg := config.Load()

	logger.InitSlog(logger.Config{
		Level:   logger.InfoLevel,
		Handler: logger.Text,
	})

	c := cache.NewInMemCache(5*time.Second, cache.EvictLRU)

	gst := server.WithGracefulShutdownTimeout(5 * time.Second)
	s := server.NewServer(":"+cfg.Port, gst)

	if cfg.ENV == "dev" {
		s.SetupSwagger()
	}
	api.SetupHealthCheck(s)

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(middleware.CORS(middleware.CORSConfig{Origin: cfg.CorsOrigins}))

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
	st := storage.NewWeatherInMemStorage()
	wService := service.NewWeatherService(wCl, astroCl, st)

	api.BindWeatherApi(s, wService, authService, c)

	s.SetupNotFoundHandler()

	go func() {
		<-s.ShutdownSig
		slog.Info("Shutdown started, cleaning up resources...")
		c.Stop()
	}()

	if err := s.Start(); err != nil {
		slog.Error(err.Error())
	}
}
