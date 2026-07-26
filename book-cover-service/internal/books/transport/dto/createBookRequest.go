package dto

type CreateBookRequest struct {
	Title       string `json:"title" example:"Moveable feast pencil drawn cover"`
	Description string `json:"description" example:"Beauty hand painted cover for my lovely book of Ernest Hemingway"`
}
