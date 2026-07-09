package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FavoritesRepository struct {
	connection *pgx.Conn
}

func NewFavoritesRepository(conn *pgx.Conn) FavoritesRepository {
	return FavoritesRepository{
		connection: conn,
	}
}

func (f *FavoritesRepository) AddCoverToFavorites(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error {
	query := `
	INSERT INTO favorites (user_id, cover_id)
	VALUES ($1, $2);
`
	_, err := f.connection.Exec(ctx, query, userID, coverID)
	if err != nil {
		return fmt.Errorf("cover repository / add cover to favorites: %w", err)
	}

	return nil
}

func (f *FavoritesRepository) GetFavoriteCoversIDs(ctx context.Context, userID uuid.UUID, offset int, limit int) ([]uuid.UUID, error) {
	query := `
	SELECT (cover_id) FROM favorites
	WHERE user_id = $1
	OFFSET $2 
	LIMIT $3;
`
	resultRow, err := f.connection.Query(ctx, query, userID, offset, limit)
	if err != nil {
		return []uuid.UUID{}, fmt.Errorf("cover repository / get favorites ids / query: %w", err)
	}

	defer resultRow.Close()

	var coversIDs []uuid.UUID

	for resultRow.Next() {
		var coverID uuid.UUID
		if err = resultRow.Scan(&coverID); err != nil {
			return []uuid.UUID{}, fmt.Errorf("cover repository / get favorites ids / parsing: %w", err)
		}
		coversIDs = append(coversIDs, coverID)
	}

	if resultRow.Err() != nil {
		return []uuid.UUID{}, fmt.Errorf("cover repository / get favorites ids / query result: %w", err)
	}

	return coversIDs, nil
}
