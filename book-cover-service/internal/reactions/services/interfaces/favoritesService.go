package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type FavoritesService interface {
	AddCoverToFavorites(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error
	GetFavoriteCovers(ctx context.Context, userID uuid.UUID, offset int, limit int) ([]domains.Cover, error)
}
