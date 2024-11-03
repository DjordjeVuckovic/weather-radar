package api

import (
	"encoding/json"
	"github.com/DjordjeVuckovic/weather-radar/internal/cache"
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/internal/model"
	"github.com/DjordjeVuckovic/weather-radar/internal/service"
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"log/slog"
	"net/http"
	"time"
)

const (
	DefaultCacheTTL = 10 * time.Minute
)

type WeatherApi struct {
	server         *server.Server
	cache          cache.Cache
	weatherService *service.WeatherService
	authService    *service.AuthService
}

func BindWeatherApi(
	s *server.Server,
	wService *service.WeatherService,
	authService *service.AuthService,
	c cache.Cache) {

	api := &WeatherApi{
		server:         s,
		weatherService: wService,
		authService:    authService,
		cache:          c,
	}

	s.GET("/api/v1/weather", api.handleWeatherByCity)
	s.POST("api/v1/weather/feedback", api.handleWeatherFeedback)
}

// handleWeatherByCity retrieves weather information for a specified city.
// @Summary Get weather by city
// @Description Get weather data for a specific city.
// @Tags weather
// @Param city query string true "City name"
// @Produce json
// @Success 200 {object} model.Weather
// @Failure 400 {object} result.Err "Validation error"
// @Failure 404 {object} result.Err "City not found"
// @Failure 500 {object} result.Err "Internal server error"
// @Failure 504 {object} result.Err "Request Timeout"
// @Router /api/v1/weather [get]
func (api *WeatherApi) handleWeatherByCity(w http.ResponseWriter, r *http.Request) error {
	city := r.URL.Query().Get("city")
	if city == "" {
		return result.ValidationErr("City query param is required")
	}

	weather, cacheHit := api.getWeatherFromCache(city)
	if cacheHit {
		return resp.WriteJSON(w, http.StatusOK, weather)
	}

	weather, err := api.weatherService.GetWeatherByCity(r.Context(), city)

	if err != nil {
		return err
	}

	api.setWeatherToCache(city, weather)

	return resp.WriteJSON(w, http.StatusOK, weather)
}

// handleWeatherFeedback handles feedback submission for weather.
// @Summary Submit weather feedback
// @Description Submit feedback about the weather in a specific city.
// @Tags weather
// @Accept json
// @Produce json
// @Param feedback body dto.WeatherFeedbackReq true "Weather feedback"
// @Success 200 {object} dto.WeatherFeedbackResp
// @Failure 400 {object} result.Err "Invalid request data"
// @Failure 401 {object} result.Err "Unauthorized"
// @Router /v1/weather/feedback [post]
// @Security BasicAuth
func (api *WeatherApi) handleWeatherFeedback(w http.ResponseWriter, r *http.Request) error {
	username, password, ok := r.BasicAuth()
	if !ok || api.authService.ValidateBasicAuth(service.AuthCredentials{
		Username: username,
		Password: password,
	}) {
		return result.UnauthorizedErr("Invalid credentials")
	}

	var feedback dto.WeatherFeedbackReq
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		return result.ValidationErr("Invalid request data")
	}

	if err := api.weatherService.SubmitFeedback(feedback); err != nil {
		return err
	}

	response := dto.WeatherFeedbackResp{
		Message: "Feedback submitted successfully",
	}
	return resp.WriteJSON(w, http.StatusOK, response)
}

func (api *WeatherApi) getWeatherFromCache(city string) (*model.Weather, bool) {
	cItem, cExist := api.cache.Get(buildCacheKey(city))
	if !cExist {
		slog.Debug("Cache miss for city", slog.String("city", city))
		return nil, false
	}
	var weather model.Weather
	err := json.Unmarshal(cItem, &weather)
	if err != nil {
		slog.Error("Failed to unmarshal weather data", slog.String("city", city))
		return nil, false
	}
	slog.Debug("Cache hit for city", slog.String("city", city))
	return &weather, true
}

func (api *WeatherApi) setWeatherToCache(city string, weather *model.Weather) {
	jsonBytes, err := json.Marshal(weather)
	if err != nil {
		slog.Error("Failed to marshal weather data", slog.String("city", city))
		return
	}
	api.cache.Set(buildCacheKey(city), jsonBytes, DefaultCacheTTL)
}

func buildCacheKey(city string) string {
	return "weather:" + city
}
