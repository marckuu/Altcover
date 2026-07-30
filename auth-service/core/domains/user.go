package domains

import (
	"auth-service/core/enums"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Nickname     string
	Role         enums.Role // Это константа из roles
	PasswordHash []byte
	CreatedAt    time.Time
}
