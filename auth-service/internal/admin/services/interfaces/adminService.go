package interfaces

import (
	"auth-service/core/enums"
	"context"

	"github.com/google/uuid"
)

type AdminService interface {
	ChangeRole(ctx context.Context, userID uuid.UUID, userRole enums.Role) error
	GetRole(ctx context.Context, userID uuid.UUID) (enums.Role, error)
}
