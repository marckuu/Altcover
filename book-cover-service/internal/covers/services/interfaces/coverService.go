package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type CoverService interface {
	GetCoverByID(ctx context.Context, coverID uuid.UUID) (domains.Cover, error)
	GetCoversByUserID(ctx context.Context, offset int, limit int, userID uuid.UUID) ([]domains.Cover, error)
	GetCoversByBook(ctx context.Context, bookID uuid.UUID, offset int, limit int) ([]domains.Cover, error)
	AddCover(ctx context.Context, cover domains.Cover) error
	UpdateCover(ctx context.Context, cover domains.Cover) error
	GetCoversByIDs(ctx context.Context, coversIDs []uuid.UUID) ([]domains.Cover, error)
	GetMostLikedCovers(ctx context.Context, daysNumber int, offset int, limit int) ([]domains.Cover, error)
}
