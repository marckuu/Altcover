package repositories

import (
	"book-cover-service/domains"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errCoverNotFound = errors.New("обложка не найдена")

type CoverRepository struct {
	connection *pgx.Conn
}

func NewCoverRepository(conn *pgx.Conn) CoverRepository {
	return CoverRepository{
		connection: conn,
	}
}

// CRUD операции

func (c *CoverRepository) AddCover(ctx context.Context, cover domains.Cover) error {
	query := `
	INSERT INTO cover (title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`
	_, err := c.connection.Exec(ctx, query,
		cover.Title,
		cover.Description,
		cover.ImagesKeys,
		cover.Status,
		cover.UserID,
		cover.DesignerNickname,
		cover.DesignerAvatarKey,
		cover.BookID)

	if err != nil {
		return fmt.Errorf("cover repository -> add: %w", err)
	}

	return nil
}

func (c *CoverRepository) GetCoverByID(ctx context.Context, coverID uuid.UUID) (domains.Cover, error) {
	query := `
	SELECT id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id
	FROM cover
	WHERE id = $1;
`
	resultRow := c.connection.QueryRow(ctx, query, coverID)

	var cover domains.Cover

	err := resultRow.Scan(&cover.ID,
		&cover.Title,
		&cover.Description,
		&cover.ImagesKeys,
		&cover.Status,
		&cover.UserID,
		&cover.DesignerNickname,
		&cover.DesignerAvatarKey,
		&cover.BookID)

	if err != nil {
		return domains.Cover{}, fmt.Errorf("cover repository -> get by id: %w", err)
	}

	return cover, nil
}

func (c *CoverRepository) GetCovers(ctx context.Context, offset int, limit int) ([]domains.Cover, error) {
	query := `
	SELECT id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id
	FROM cover
	OFFSET $1
	LIMIT $2;
`
	resultRows, err := c.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var covers []domains.Cover

	for resultRows.Next() {
		cover := domains.Cover{}

		err = resultRows.Scan(&cover.ID,
			&cover.Title,
			&cover.Description,
			&cover.ImagesKeys,
			&cover.Status,
			&cover.UserID,
			&cover.DesignerNickname,
			&cover.DesignerAvatarKey,
			&cover.BookID)

		if err != nil {
			return []domains.Cover{}, fmt.Errorf("cover repository -> get all -> parsing: %w", err)
		}

		covers = append(covers, cover)
	}

	if resultRows.Err() != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get all -> query result: %w", err)
	}

	return covers, nil
}

func (c *CoverRepository) UpdateCover(ctx context.Context, cover domains.Cover) error {
	query := `
	UPDATE cover
	SET title = $1, 
	    description = $2,  
	    images_keys = $3, 
	    status = $4, 
	    designer_nickname = $5, 
	    designer_avatar_key = $6,
	    book_id = $7
	WHERE id = $8;
`
	tag, err := c.connection.Exec(ctx, query,
		cover.Title,
		cover.Description,
		cover.ImagesKeys,
		cover.Status,
		cover.DesignerNickname,
		cover.DesignerAvatarKey,
		cover.BookID,
		cover.ID)

	if err != nil {
		return fmt.Errorf("cover repository -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errCoverNotFound
	}

	return nil
}

func (c *CoverRepository) DeleteCover(ctx context.Context, coverID uuid.UUID) error {
	query := `
	DELETE FROM cover
	WHERE id = $1;
`
	tag, err := c.connection.Exec(ctx, query, coverID)
	if err != nil {
		return fmt.Errorf("cover repository -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errCoverNotFound
	}

	return nil
}

// Дополнительные операции

func (c *CoverRepository) GetCoversByUserID(ctx context.Context, offset int, limit int, userID uuid.UUID) ([]domains.Cover, error) {
	query := `
	SELECT id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id
	FROM cover
	WHERE user_id = $1
	OFFSET $2
	LIMIT $3;
`
	resultRows, err := c.connection.Query(ctx, query, userID, offset, limit)
	if err != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get all by user id -> query: %w", err)
	}

	defer resultRows.Close()

	var covers []domains.Cover

	for resultRows.Next() {
		var cover domains.Cover

		err = resultRows.Scan(&cover.ID,
			&cover.Title,
			&cover.Description,
			&cover.ImagesKeys,
			&cover.Status,
			&cover.UserID,
			&cover.DesignerNickname,
			&cover.DesignerAvatarKey,
			&cover.BookID)

		if err != nil {
			return []domains.Cover{}, fmt.Errorf("cover repository -> get all by user id -> parsing: %w", err)
		}

		covers = append(covers, cover)
	}

	if resultRows.Err() != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get all by user id -> query result: %w", err)
	}

	return covers, nil
}

func (c *CoverRepository) GetCoversByIDs(ctx context.Context, coversIDs []uuid.UUID) ([]domains.Cover, error) {
	query := `
	SELECT id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id
	FROM cover
	WHERE id = ANY($1);
`
	resultRow, err := c.connection.Query(ctx, query, coversIDs)
	if err != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get covers by ids -> query: %w", err)
	}

	defer resultRow.Close()

	var covers []domains.Cover

	for resultRow.Next() {
		var cover domains.Cover
		err = resultRow.Scan(&cover.ID,
			&cover.Title,
			&cover.Description,
			&cover.ImagesKeys,
			&cover.Status,
			&cover.UserID,
			&cover.DesignerNickname,
			&cover.DesignerAvatarKey,
			&cover.BookID)

		if err != nil {
			return []domains.Cover{}, fmt.Errorf("cover repository -> get covers by ids -> parsing: %w", err)
		}

		covers = append(covers, cover)
	}

	if resultRow.Err() != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get covers by ids -> query result: %w", err)
	}

	return covers, nil
}

func (c *CoverRepository) GetMostLikedCoversForNDays(ctx context.Context, daysNumber int, offset int, limit int) ([]domains.Cover, error) {
	query := `
	SELECT (id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id), COUNT(cover_like.user_id) AS likes_count
	FROM cover c
	INNER JOIN cover_like l
	ON c.id=l.user_id
	AND l.created_at >= now() - ($1 *interval '1 day')
	GROUP BY c.id
	ORDER BY likes_count DESC
	OFFSET $2
	LIMIT $3;
`
	resultRow, err := c.connection.Query(ctx, query, daysNumber, offset, limit)
	if err != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get most liked -> query: %w", err)
	}

	defer resultRow.Close()

	var covers []domains.Cover

	for resultRow.Next() {
		var cover domains.Cover

		err = resultRow.Scan(&cover.ID,
			&cover.Title,
			&cover.Description,
			&cover.ImagesKeys,
			&cover.Status,
			&cover.UserID,
			&cover.DesignerNickname,
			&cover.DesignerAvatarKey,
			&cover.BookID)

		if err != nil {
			return []domains.Cover{}, fmt.Errorf("cover repository -> get most liked -> parsing: %w", err)
		}

		covers = append(covers, cover)
	}

	if resultRow.Err() != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get most liked -> query result: %w", err)
	}

	return covers, nil
}

func (c *CoverRepository) GetCoversByBook(ctx context.Context, bookID uuid.UUID, offset int, limit int) ([]domains.Cover, error) {
	query := `
	SELECT (id, title, description, images_keys, status, user_id, designer_nickname, designer_avatar_key, book_id)
	FROM cover
	WHERE book_id = $1
	OFFSET $2
	LIMIT $3;
`
	resultRow, err := c.connection.Query(ctx, query, bookID, offset, limit)
	if err != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get most liked -> query: %w", err)
	}

	defer resultRow.Close()

	var covers []domains.Cover

	for resultRow.Next() {
		var cover domains.Cover

		err = resultRow.Scan(&cover.ID,
			&cover.Title,
			&cover.Description,
			&cover.ImagesKeys,
			&cover.Status,
			&cover.UserID,
			&cover.DesignerNickname,
			&cover.DesignerAvatarKey,
			&cover.BookID)

		if err != nil {
			return []domains.Cover{}, fmt.Errorf("cover repository -> get by book -> parsing: %w", err)
		}

		covers = append(covers, cover)
	}

	if resultRow.Err() != nil {
		return []domains.Cover{}, fmt.Errorf("cover repository -> get by book -> query result: %w", err)
	}

	return covers, nil
}
