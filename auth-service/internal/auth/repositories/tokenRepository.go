package repositories

import (
	"auth-service/core/db"
	"auth-service/core/domains"
	errors2 "auth-service/core/errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TokenRepository struct {
	connection db.Database
}

func NewTokenRepository(conn db.Database) *TokenRepository {
	return &TokenRepository{
		connection: conn,
	}
}

func (t *TokenRepository) AddRefreshToken(ctx context.Context, tokenHash []byte) error {
	query := `
	INSERT INTO refresh_token (token_hash)
	VALUES ($1);
`
	_, err := t.connection.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("tokens repo -> add: %w", err)
	}

	return nil
}

func (t *TokenRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash []byte) (domains.Token, error) {
	query := `
	SELECT id, token_hash
	FROM refresh_token
	WHERE token_hash = $1;
`
	resultRow := t.connection.QueryRow(ctx, query, tokenHash)

	var token domains.Token

	if err := resultRow.Scan(&token.ID, &token.TokenHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Token{}, errors2.ErrTokenNotFound
		}
		return domains.Token{}, fmt.Errorf("tokens repo -> get by hash: %w", err)
	}

	return token, nil
}

func (t *TokenRepository) DeleteRefreshToken(ctx context.Context, tokenHash []byte) error {
	query := `
	DELETE 
	FROM refresh_token
	WHERE token_hash = $1;
`
	tag, err := t.connection.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("tokens repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errors2.ErrTokenNotFound
	}

	return nil
}
