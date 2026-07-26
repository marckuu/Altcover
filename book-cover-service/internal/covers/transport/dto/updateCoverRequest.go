package dto

type UpdateCoverRequest struct {
	Title       string   `json:"title" example:"Moveable fest cover"`
	Description string   `json:"description" example:"Ernest Hemingway Moveable fest anime style cover"`
	ImagesKeys  []string `json:"imagesKeys" example:"/covers/1252"`
	Status      int      `json:"status" example:"0"`
}
