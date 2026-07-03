package domains

import (
	"github.com/google/uuid"
)

type DesignerProfile struct {
	ID        uuid.UUID `json:"id"`
	Nickname  string    `json:"nickname"`
	AvatarKey string    `json:"avatarKey"` // Нужно будет настроить работу с minio

	UserID uuid.UUID `json:"userID"`
}
