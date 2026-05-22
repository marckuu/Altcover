package dto

import "time"

type ErrorResponse struct {
	Message string
	Time    time.Time
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Time:    time.Now(),
	}
}
