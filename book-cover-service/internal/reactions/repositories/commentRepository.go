package repositories

import (
	"book-cover-service/core/domains"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var errCommentNotFound = errors.New("комментарий не найден")

type CommentRepository struct {
	connection *pgx.Conn
}

func (c *CommentRepository) AddComment(ctx context.Context, comment domains.Comment) error {
	query := `
	INSERT INTO comment (text, cover_id, user_id)
	VALUES ($1, $2, $3);
`
	if _, err := c.connection.Exec(ctx, query, comment.Text, comment.CoverId, comment.UserId); err != nil {
		return fmt.Errorf("comment repo -> add: %w", err)
	}

	return nil
}

func (c *CommentRepository) GetCommentByID(ctx context.Context, commentID int64) (domains.Comment, error) {
	query := `
	SELECT id, text, cover_id, user_id
	FROM comment
	WHERE id = $1;
`
	resultRow := c.connection.QueryRow(ctx, query, commentID)

	var comment domains.Comment

	if err := resultRow.Scan(&comment.ID, &comment.Text, &comment.CoverId); err != nil {
		return domains.Comment{}, fmt.Errorf("comment repo -> get by id: %w", err)
	}

	return comment, nil
}

func (c *CommentRepository) GetComments(ctx context.Context, offset int, limit int) ([]domains.Comment, error) {
	query := `
	SELECT id, text, cover_id, user_id
	FROM comment
	OFFSET $1
	LIMIT $2
`
	resultRows, err := c.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains.Comment{}, fmt.Errorf("comment repo -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var comments []domains.Comment

	for resultRows.Next() {
		var comment domains.Comment

		if err = resultRows.Scan(&comment.ID, &comment.Text, &comment.CoverId); err != nil {
			return []domains.Comment{}, fmt.Errorf("comment repo -> get all -> parsing: %w", err)
		}

		comments = append(comments, comment)
	}

	if resultRows.Err() != nil {
		return []domains.Comment{}, fmt.Errorf("comment repo -> get all -> query result: %w", err)
	}

	return comments, nil
}

func (c *CommentRepository) UpdateComment(ctx context.Context, comment domains.Comment) error {
	query := `
	UPDATE comment
	SET text = $1, cover_id = $2, user_id = $3
	WHERE id = $3;
`
	tag, err := c.connection.Exec(ctx, query, comment.Text, comment.CoverId, comment.UserId, comment.ID)
	if err != nil {
		return fmt.Errorf("comment repo -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errCommentNotFound
	}

	return nil
}

func (c *CommentRepository) DeleteComment(ctx context.Context, commentID int64) error {
	query := `
	DELETE FROM comment
	WHERE id = $1;
`
	tag, err := c.connection.Exec(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("comment repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errCommentNotFound
	}

	return nil
}
