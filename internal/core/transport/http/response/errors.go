package core_http_response

type ErrorResponse struct {
	Error   string `json:"error" example:"bad request"`
	Message string `json:"message" example:"потеря потерь"`
}
