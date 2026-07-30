package interfaces

import (
	"auth-service/core/enums"
	"context"

	"github.com/google/uuid"
)

type AdminRepository interface {
	UpdateRole(ctx context.Context, userID uuid.UUID, userRole enums.Role) error
	GetRole(ctx context.Context, userID uuid.UUID) (enums.Role, error)
}
