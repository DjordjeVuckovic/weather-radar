package resp

import (
	"encoding/json"
	results "github.com/DjordjeVuckovic/weather-radar/pkg/result"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, code int, body interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}

func WriteProblemJSON(w http.ResponseWriter, p *results.Err) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	return json.NewEncoder(w).Encode(p)
}
