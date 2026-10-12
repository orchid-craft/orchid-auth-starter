package model

type Response struct {
	Status       int    `json:"status"`
	Success      bool   `json:"success"`
	Data         any    `json:"data,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
