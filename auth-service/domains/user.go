package domains

import "time"

type User struct {
	ID           int64
	Nickname     string
	Role         int // Это константа из roles
	PasswordHash string
	CreatedAt    time.Time
}
