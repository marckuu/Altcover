package dto

import "github.com/google/uuid"

type DesignerProfileSnapshot struct {
	ID        uuid.UUID `json:"id"`
	AvatarKey string    `json:"avatarKey"`
	Nickname  string    `json:"nickname"`
}
