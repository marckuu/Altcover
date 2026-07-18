package repositories

import (
	"auth-service/core/db"
	"auth-service/core/domains"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errDesignerProfileNotFound = errors.New("профиль дизайнера не найден")

type DesignerProfileRepository struct {
	connection db.Database
}

func NewDesignerProfileRepository(conn *pgx.Conn) *DesignerProfileRepository {
	return &DesignerProfileRepository{
		connection: conn,
	}
}

func (d *DesignerProfileRepository) AddDesignerProfile(ctx context.Context, profile domains.DesignerProfile) error {
	query := `
	INSERT INTO designer_profile (user_id, avatar_key, nickname)
	VALUES ($1, $2, $3);
`
	if _, err := d.connection.Exec(ctx, query, profile.UserID, profile.AvatarKey, profile.Nickname); err != nil {
		return fmt.Errorf("designer profile repo -> add: %w", err)
	}

	return nil
}

func (d *DesignerProfileRepository) GetDesignerProfileByID(ctx context.Context, profileID int64) (domains.DesignerProfile, error) {
	query := `
	SELECT id, user_id, avatar_key, nickname
	FROM designer_profile
	WHERE id = $1;
`
	resultRow := d.connection.QueryRow(ctx, query, profileID)

	var designerProfile domains.DesignerProfile

	if err := resultRow.Scan(&designerProfile.ID, &designerProfile.UserID, &designerProfile.AvatarKey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.DesignerProfile{}, errDesignerProfileNotFound
		}
		return domains.DesignerProfile{}, fmt.Errorf("designer profile repo -> get by id: %w", err)
	}

	return designerProfile, nil
}

func (d *DesignerProfileRepository) GetDesignersProfiles(ctx context.Context, offset int, limit int) ([]domains.DesignerProfile, error) {
	query := `
	SELECT id, user_id, avatar_key, nickname
	FROM designer_profile
	OFFSET $1
	LIMIT $2;
`
	resultRows, err := d.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains.DesignerProfile{}, fmt.Errorf("designer profile repo -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var profiles []domains.DesignerProfile

	for resultRows.Next() {
		profile := domains.DesignerProfile{}
		if err = resultRows.Scan(&profile.ID, &profile.UserID, &profile.AvatarKey); err != nil {
			return []domains.DesignerProfile{}, fmt.Errorf("designer profile repo -> get all -> parsing: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err = resultRows.Err(); err != nil {
		return []domains.DesignerProfile{}, fmt.Errorf("designer profile repo -> get all -> query thread: %w", err)
	}

	return profiles, nil
}

func (d *DesignerProfileRepository) UpdateDesignerProfile(ctx context.Context, profile domains.DesignerProfile) error {
	query := `
	UPDATE designer_profile
	SET avatar_key = $1, 
		nickname = $2
	WHERE id = $3;
`
	tag, err := d.connection.Exec(ctx, query, profile.AvatarKey, profile.Nickname, profile.ID)
	if err != nil {
		return fmt.Errorf("designer profile repo -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errDesignerProfileNotFound
	}

	return nil
}

func (d *DesignerProfileRepository) DeleteDesignerProfile(ctx context.Context, profileID uuid.UUID) error {
	query := `
	DELETE FROM designer_profile
	WHERE id = $1;
`
	tag, err := d.connection.Exec(ctx, query, profileID)
	if err != nil {
		return fmt.Errorf("designer profile repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errDesignerProfileNotFound
	}

	return nil
}

func (d *DesignerProfileRepository) GetDesignerProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error) {
	query := `
	SELECT id, user_id, avatar_key, nickname
	FROM designer_profile
	WHERE user_id = $1;
`
	resultRow := d.connection.QueryRow(ctx, query, userID)

	var designerProfile domains.DesignerProfile

	if err := resultRow.Scan(&designerProfile.ID, &designerProfile.UserID, &designerProfile.AvatarKey, &designerProfile.Nickname); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.DesignerProfile{}, errDesignerProfileNotFound
		}
		return domains.DesignerProfile{}, fmt.Errorf("designer profile repo -> get by user id: %w", err)
	}

	return designerProfile, nil
}
