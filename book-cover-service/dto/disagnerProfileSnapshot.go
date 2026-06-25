package dto

import "github.com/google/uuid"

type DesignerProfileSnapshot struct {
	ID        uuid.UUID
	avatarKey string
	Nickname  string
}
