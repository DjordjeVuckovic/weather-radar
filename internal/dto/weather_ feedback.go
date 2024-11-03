package dto

type WeatherFeedbackResp struct {
	Message string
}

type WeatherFeedbackReq struct {
	Date    string `json:"date"`
	City    string `json:"city"`
	Message string `json:"Message"`
}
