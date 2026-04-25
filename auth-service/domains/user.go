package domains

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Nickname     string
	Role         int // Это константа из roles
	PasswordHash string
	CreatedAt    time.Time
}
