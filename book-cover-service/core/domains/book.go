package domains

import (
	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Description string
}
