package repositories

import (
	"book-cover-service/core/domains"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errLikeNotFound = errors.New("лайк не найден")

type CoverLikeRepository struct {
	connection *pgx.Conn
}

func NewCoverLikeRepository(conn *pgx.Conn) CoverLikeRepository {
	return CoverLikeRepository{
		connection: conn,
	}
}

func (lr *CoverLikeRepository) AddLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error {
	query := `
	INSERT INTO cover_like (user_id, cover_id)
	VALUES ($1, $2);
`
	if _, err := lr.connection.Exec(ctx, query, userID, coverID); err != nil {
		return fmt.Errorf("like repo -> add: %w", err)
	}

	return nil
}

func (lr *CoverLikeRepository) GetLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) (domains.CoverLike, error) {
	query := `
	SELECT (user_id, cover_id, created_at)
	FROM cover_like
	WHERE user_id = $1 AND cover_id = $2;
`
	resultRow := lr.connection.QueryRow(ctx, query, userID, coverID)

	var coverLike domains.CoverLike

	if err := resultRow.Scan(&coverLike.UserID, &coverLike.CoverID, &coverLike.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.CoverLike{}, errLikeNotFound
		}
		return domains.CoverLike{}, err
	}

	return coverLike, nil
}
