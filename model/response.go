package model

type Response struct {
	Status  any    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}
