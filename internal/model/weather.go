package model

import (
	"github.com/DjordjeVuckovic/weather-radar/internal/dto"
	"github.com/DjordjeVuckovic/weather-radar/pkg/util"
)

type Location struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	TzId      string  `json:"tz_id"`
	Localtime string  `json:"localtime"`
	TzOffset  int     `json:"tz_offset"`
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

func NewWeatherFromDto(weatherDto *dto.WeatherByCity, astroDto *dto.AstroByCity) *Weather {
	location := Location{
		Name:      weatherDto.Location.Name,
		Region:    weatherDto.Location.Region,
		Country:   weatherDto.Location.Country,
		Lat:       weatherDto.Location.Lat,
		Lon:       weatherDto.Location.Lon,
		TzId:      weatherDto.Location.TzId,
		Localtime: weatherDto.Location.Localtime,
		TzOffset:  astroDto.Timezone,
	}

	current := Current{
		LastUpdated: weatherDto.Current.LastUpdated,
		TempC:       weatherDto.Current.TempC,
		Condition:   weatherDto.Current.Condition.Text,
		WindKph:     weatherDto.Current.WindKph,
		WindDegree:  weatherDto.Current.WindDegree,
		WindDir:     weatherDto.Current.WindDir,
		PressureMb:  int(weatherDto.Current.PressureMb),
		PrecipMm:    int(weatherDto.Current.PrecipMm),
		Humidity:    weatherDto.Current.Humidity,
		Cloud:       weatherDto.Current.Cloud,
		FeelslikeC:  weatherDto.Current.FeelslikeC,
		HeatindexC:  weatherDto.Current.HeatindexC,
		VisKm:       int(weatherDto.Current.VisKm),
		Uv:          int(weatherDto.Current.Uv),
	}

	sunrise := util.UnixToLocal(int64(astroDto.Sys.Sunrise), astroDto.Timezone)
	sunset := util.UnixToLocal(int64(astroDto.Sys.Sunset), astroDto.Timezone)

	astroData := Astro{
		Sunrise: sunrise,
		Sunset:  sunset,
	}

	return &Weather{
		Location: location,
		Current:  current,
		Astro:    astroData,
	}
}
