package interfaces

import (
	"auth-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileRepository interface {
	GetDesignerProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error)
	DeleteDesignerProfile(ctx context.Context, profileID uuid.UUID) error
	UpdateDesignerProfile(ctx context.Context, profile domains.DesignerProfile) error
	AddDesignerProfile(ctx context.Context, profile domains.DesignerProfile) error
}
