package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type CoverRepository interface {
	AddCover(ctx context.Context, cover domains.Cover) (domains.Cover, error)
	GetCoverByID(ctx context.Context, coverID uuid.UUID) (domains.Cover, error)
	GetCoversByUserID(ctx context.Context, offset int, limit int, userID uuid.UUID) ([]domains.Cover, error)
	GetCoversByIDs(ctx context.Context, coversIDs []uuid.UUID) ([]domains.Cover, error)
	GetCoversByBook(ctx context.Context, bookID uuid.UUID, offset int, limit int) ([]domains.Cover, error)
	GetCovers(ctx context.Context, offset int, limit int) ([]domains.Cover, error)
	UpdateCover(ctx context.Context, cover domains.Cover) (domains.Cover, error)
	DeleteCover(ctx context.Context, coverID uuid.UUID) error
	GetMostLikedCoversForNDays(ctx context.Context, daysNumber int, offset int, limit int) ([]domains.Cover, error)
}
