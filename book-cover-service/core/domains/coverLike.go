package domains

import (
	"time"

	"github.com/google/uuid"
)

type CoverLike struct {
	UserID    uuid.UUID
	CoverID   uuid.UUID
	CreatedAt time.Time
}
