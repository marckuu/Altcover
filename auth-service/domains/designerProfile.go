package domains

import (
	"github.com/google/uuid"
)

type DesignerProfile struct {
	ID        uuid.UUID
	Nickname  string
	AvatarKey string // Нужно будет настроить работу с minio

	UserID uuid.UUID
}
