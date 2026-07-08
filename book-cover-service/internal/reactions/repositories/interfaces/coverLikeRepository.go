package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type CoverLikeRepository interface {
	AddLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error
	GetLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) (domains.CoverLike, error)
}
