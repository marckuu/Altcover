package domains

import (
	"github.com/google/uuid"
)

type Comment struct {
	ID   uuid.UUID
	Text string

	CoverId uuid.UUID
	UserId  uuid.UUID
}
