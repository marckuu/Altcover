package domains

import (
	"github.com/google/uuid"
)

type Favorites struct {
	userID  uuid.UUID
	coverID uuid.UUID
}
