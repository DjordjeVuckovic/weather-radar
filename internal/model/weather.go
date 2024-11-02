package model

type Location struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	TzId      string  `json:"tz_id"`
	Localtime string  `json:"localtime"`
}

type Current struct {
	LastUpdated string  `json:"last_updated"`
	TempC       float64 `json:"temp_c"`
	Condition   string  `json:"condition"`
	WindKph     float64 `json:"wind_kph"`
	WindDegree  int     `json:"wind_degree"`
	WindDir     string  `json:"wind_dir"`
	PressureMb  int     `json:"pressure_mb"`
	PrecipMm    int     `json:"precip_mm"`
	Humidity    int     `json:"humidity"`
	Cloud       int     `json:"cloud"`
	FeelslikeC  float64 `json:"feelslike_c"`
	HeatindexC  float64 `json:"heatindex_c"`
	VisKm       int     `json:"vis_km"`
	Uv          int     `json:"uv"`
}
type Astro struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type Weather struct {
	Location `json:"location"`
	Current  `json:"current"`
	Astro    `json:"astro"`
}
