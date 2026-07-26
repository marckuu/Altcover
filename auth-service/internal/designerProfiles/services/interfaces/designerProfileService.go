package interfaces

import (
	"auth-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileService interface {
	CreateDesignerProfileToUser(ctx context.Context, userID uuid.UUID, designerProfile domains.DesignerProfile) (domains.DesignerProfile, error)
	UpdateDesignerProfileToUser(ctx context.Context, userID uuid.UUID, newDesignerProfile domains.DesignerProfile) (domains.DesignerProfile, error)
	DeleteDesignerProfileToUser(ctx context.Context, userID uuid.UUID) error
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error)
}
