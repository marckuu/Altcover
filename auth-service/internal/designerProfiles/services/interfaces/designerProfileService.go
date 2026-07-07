package interfaces

import (
	"auth-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileService interface {
	DeleteDesignerProfile(ctx context.Context, profileID uuid.UUID) error
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error)
	UpdateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error
	CreateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error
}
