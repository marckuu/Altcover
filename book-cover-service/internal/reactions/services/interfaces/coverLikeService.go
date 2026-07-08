package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type CoverLikeService interface {
	SetLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error
}
