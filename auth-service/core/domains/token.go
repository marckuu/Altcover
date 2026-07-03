package domains

import "github.com/google/uuid"

type Token struct {
	ID        uuid.UUID
	TokenHash []byte
}
