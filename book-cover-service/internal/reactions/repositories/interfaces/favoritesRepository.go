package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type FavoritesRepository interface {
	AddCoverToFavorites(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error
	GetFavoriteCoversIDs(ctx context.Context, userID uuid.UUID, offset int, limit int) ([]uuid.UUID, error)
}
