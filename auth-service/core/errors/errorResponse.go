package errors

import "time"

type ErrorResponse struct {
	Message string    `json:"message" example:"some error..."`
	Time    time.Time `json:"time" example:"2026-04-12:13:20:47"`
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Time:    time.Now(),
	}
}
