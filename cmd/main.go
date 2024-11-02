package main

import (
	api2 "github.com/DjordjeVuckovic/weather-radar/api"
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/env"
	"github.com/DjordjeVuckovic/weather-radar/pkg/logger"
	"github.com/DjordjeVuckovic/weather-radar/pkg/middleware"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"os"
	"time"
)

func main() {
	env.Load()

	logger.InitSlog(logger.Config{
		Level:   logger.InfoLevel,
		Handler: logger.Text,
	})

	gst := server.WithGracefulShutdownTimeout(5 * time.Second)
	s := server.NewServer(":1312", gst)

	s.Use(middleware.Logger())
	s.Use(middleware.Recover())
	s.Use(middleware.CORS(middleware.Config{Origin: "*"}))
	s.Use(middleware.Example())

	wCl := client.NewAPIWeatherClient(os.Getenv("WEATHER_API_URL"), os.Getenv("WEATHER_API_KEY"))
	api2.BindWeatherApi(s, wCl)

	s.SetupNotFoundHandler()

	if err := s.Start(); err != nil {
		slog.Error(err.Error())
	}
}
