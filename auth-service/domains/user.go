package domains

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           pgtype.UUID
	Nickname     string
	Role         int // Это константа из roles
	PasswordHash string
	CreatedAt    time.Time
}
