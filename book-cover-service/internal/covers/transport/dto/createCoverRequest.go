package dto

import "github.com/google/uuid"

type CreateCoverRequest struct {
	Title       string   `json:"title" example:"Moveable fest cover"`
	Description string   `json:"description" example:"Ernest Hemingway Moveable fest anime style cover"`
	ImagesKeys  []string `json:"imagesKeys" example:"/covers/1252"`
	Status      int      `json:"status" example:"0"`

	BookID uuid.UUID `json:"book_id" example:"-> paste book's uuid <-"`
}
