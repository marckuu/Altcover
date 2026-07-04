package repositories

import (
	"book-cover-service/internal/snapshots/transport/dto"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errDesignerProfileSnapshotNotFound = errors.New("снимок профиля дизайнера не найден")

type DesignerProfileSnapshotRepository struct {
	connection *pgx.Conn
}

func NewDesignerProfileSnapshotRepository(conn *pgx.Conn) DesignerProfileSnapshotRepository {
	return DesignerProfileSnapshotRepository{
		connection: conn,
	}
}

func (d *DesignerProfileSnapshotRepository) AddDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error {
	query := `
	INSERT INTO designer_profile_snapshot (id, avatar_key, nickname, user_id)
	VALUES ($1, $2, $3, $4);
`
	if _, err := d.connection.Exec(ctx, query, profile.ID, profile.AvatarKey, profile.Nickname, profile.UserID); err != nil {
		return fmt.Errorf("designer profile snapshot repo / add: %w", err)
	}

	return nil
}

func (d *DesignerProfileSnapshotRepository) UpdateDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error {
	query := `
	UPDATE designer_profile_snapshot
	SET avatar_key = $1, nickname = $2, user_id = $3
	WHERE id = $4;
`
	tag, err := d.connection.Exec(ctx, query, profile.AvatarKey, profile.Nickname, profile.UserID, profile.ID)
	if err != nil {
		return fmt.Errorf("designer profile snapshot repo / update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errDesignerProfileSnapshotNotFound
	}

	return nil
}

func (d *DesignerProfileSnapshotRepository) DeleteDesignerProfileSnapshot(ctx context.Context, profileID uuid.UUID) error {
	query := `
	DELETE FROM designer_profile_snapshot
	WHERE id = $1;
`
	tag, err := d.connection.Exec(ctx, query, profileID)
	if err != nil {
		return fmt.Errorf("comment repo / delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errDesignerProfileSnapshotNotFound
	}

	return nil
}

func (d *DesignerProfileSnapshotRepository) GetDesignerProfileSnapshotByUserID(ctx context.Context, userID uuid.UUID) (dto.DesignerProfileSnapshot, error) {
	query := `
	SELECT id, avatar_key, nickname, user_id
	FROM designer_profile_snapshot
	WHERE user_id = $1;
`

	resultRow := d.connection.QueryRow(ctx, query, userID)

	var profileSnapshot dto.DesignerProfileSnapshot

	if err := resultRow.Scan(&profileSnapshot.ID, &profileSnapshot.AvatarKey, &profileSnapshot.Nickname, &profileSnapshot.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.DesignerProfileSnapshot{}, errDesignerProfileSnapshotNotFound
		}
		return dto.DesignerProfileSnapshot{}, err
	}

	return profileSnapshot, nil
}
