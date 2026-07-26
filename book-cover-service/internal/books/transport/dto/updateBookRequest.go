package dto

import "github.com/google/uuid"

type UpdateBookRequest struct {
	ID          uuid.UUID `json:"id" example:"-> paste book's uuid <-"`
	Title       string    `json:"title" example:"Moveable feast anime style cover"`
	Description string    `json:"description" example:"Anime style cover for Ernest Hemingway Moveable feast"`
}
