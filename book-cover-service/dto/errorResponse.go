package dto

import "time"

type ErrorResponse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Time:    time.Now(),
	}
}
