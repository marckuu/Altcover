package repositories

import (
	"auth-service/core/domains"
	"auth-service/core/tools"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testTimeout = 10 * time.Second

func TestTokenRepository_AddRefreshToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	token := "refreshToken"
	tokenHash, err := tools.GetTokenHash(token)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	tk := &TokenRepository{
		connection: tx,
	}

	err = tk.AddRefreshToken(ctx, tokenHash)
	require.NoError(t, err)

	receivedToken := &domains.Token{
		ID:        uuid.UUID{},
		TokenHash: nil,
	}
	err = tx.QueryRow(ctx, `
	SELECT token_hash 
	FROM refresh_token  
	WHERE token_hash=$1;
`, tokenHash).Scan(&receivedToken.TokenHash)
	require.NoError(t, err)

	assert.Equal(t, tokenHash, receivedToken.TokenHash)
}

func TestTokenRepository_GetRefreshTokenByHash(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	token := "refreshToken"
	tokenHash, err := tools.GetTokenHash(token)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	tk := &TokenRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO refresh_token (token_hash)
	VALUES ($1);
`, tokenHash)
	require.NoError(t, err)

	refreshToken, err := tk.GetRefreshTokenByHash(ctx, tokenHash)
	require.NoError(t, err)

	assert.Equal(t, tokenHash, refreshToken.TokenHash)
}

func TestTokenRepository_DeleteRefreshToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	token := "refreshToken"
	tokenHash, err := tools.GetTokenHash(token)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	tk := &TokenRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO refresh_token (token_hash)
	VALUES ($1);
`, tokenHash)
	require.NoError(t, err)

	err = tk.DeleteRefreshToken(ctx, tokenHash)
	require.NoError(t, err)

	receivedToken := &domains.Token{}
	err = tx.QueryRow(ctx, `
	SELECT (token_hash)
	FROM refresh_token
	WHERE token_hash=$1;
`, tokenHash).Scan(receivedToken.TokenHash)

	assert.Equal(t, pgx.ErrNoRows, err)
}
