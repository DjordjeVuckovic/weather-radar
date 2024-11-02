package api

import (
	"github.com/DjordjeVuckovic/weather-radar/internal/client"
	"github.com/DjordjeVuckovic/weather-radar/internal/usecase"
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
)

type WeatherApi struct {
	server *server.Server
	wCl    client.WeatherClient
}

func BindWeatherApi(s *server.Server, wCl client.WeatherClient) {
	api := &WeatherApi{
		server: s,
		wCl:    wCl,
	}

	s.GET("/api/v1/weather", api.handleWeatherByCity)
}

// handleWeatherByCity @Tags weather
// @QueryParam city
// @Tags rooms
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.PokerRoomDto
// @Failure 404 {object} result.Problem
// @Router /api/v1/weather?city="" [get]
func (api *WeatherApi) handleWeatherByCity(w http.ResponseWriter, r *http.Request) error {
	city := r.URL.Query().Get("city")
	if city == "" {
		return result.ValidationErr("City query param is required")
	}

	weather, err := usecase.WeatherByCityQuery(r.Context(), api.wCl, city)

	if err != nil {
		return err
	}

	return resp.WriteJSON(w, http.StatusOK, weather)
}
