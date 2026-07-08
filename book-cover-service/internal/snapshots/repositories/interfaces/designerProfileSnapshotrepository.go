package interfaces

import (
	"book-cover-service/internal/snapshots/transport/dto"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileSnapshotRepository interface {
	AddDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error
	GetDesignerProfileSnapshotByUserID(ctx context.Context, userID uuid.UUID) (dto.DesignerProfileSnapshot, error)
	UpdateDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error
	DeleteDesignerProfileSnapshot(ctx context.Context, profileID uuid.UUID) error
}
