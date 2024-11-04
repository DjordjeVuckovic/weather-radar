package api

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/resp"
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"net/http"
)

// SetupHealthCheck sets up health and readiness checks for the server.
func SetupHealthCheck(s *server.Server) {

	s.GET("/healthz", handleHealthChecks)
	s.GET("/ready", handleReadinessChecks)
}

// @Summary Health check endpoint
// @Description This endpoint returns the health status of the application.
// @Tags health
// @Produce json
// @Success 200 {string} string "OK"
// @Router /healthz [get]
func handleHealthChecks(w http.ResponseWriter, _ *http.Request) error {
	err := resp.WriteJSON(w, http.StatusOK, "OK")
	if err != nil {
		return err
	}
	return nil
}

// @Summary Readiness check endpoint
// @Description This endpoint returns the readiness status of the application.
// @Tags health
// @Produce json
// @Success 200 {string} string "OK"
// @Router /ready [get]
func handleReadinessChecks(w http.ResponseWriter, _ *http.Request) error {
	err := resp.WriteJSON(w, http.StatusOK, "OK")
	if err != nil {
		return err
	}
	return nil
}
