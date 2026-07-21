package repositories

import (
	"auth-service/core/domains"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepository_AddUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	require.NoError(t, err)

	user := domains.User{
		Nickname:     "Ivan",
		Role:         1,
		PasswordHash: passwordHash,
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	u := &UserRepository{
		connection: tx,
	}

	err = u.AddUser(ctx, user)
	require.NoError(t, err)

	var savedUser domains.User
	err = tx.QueryRow(ctx, `
	SELECT nickname, role, password_hash
	FROM users
	WHERE nickname=$1;
`, user.Nickname).Scan(&savedUser.Nickname, &savedUser.Role, &savedUser.PasswordHash)

	require.NoError(t, err)

	assert.Equal(t, user.Nickname, savedUser.Nickname)
	assert.Equal(t, user.Role, savedUser.Role)
	assert.Equal(t, user.PasswordHash, savedUser.PasswordHash)

}

func TestUserRepository_GetUserByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	userID := uuid.New()
	nickname := "Ivan"
	role := 1
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	u := &UserRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (id, nickname, role, password_hash) VALUES ($1, $2, $3, $4);
`, userID, nickname, role, passwordHash)
	require.NoError(t, err)

	receivedUser, err := u.GetUserByID(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, userID, receivedUser.ID)
	assert.Equal(t, nickname, receivedUser.Nickname)
	assert.Equal(t, role, receivedUser.Role)
	assert.Equal(t, passwordHash, receivedUser.PasswordHash)
}

func TestUserRepository_GetUserByNickname(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	nickname := "Ivan"
	role := 1
	userID := uuid.New()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	u := &UserRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (id, nickname, role, password_hash) VALUES ($1, $2, $3, $4);
`, userID, nickname, role, passwordHash)
	require.NoError(t, err)

	receivedUser, err := u.GetUserByNickname(ctx, nickname)
	require.NoError(t, err)

	assert.Equal(t, userID, receivedUser.ID)
	assert.Equal(t, nickname, receivedUser.Nickname)
	assert.Equal(t, role, receivedUser.Role)
	assert.Equal(t, passwordHash, receivedUser.PasswordHash)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	userID := uuid.New()
	nickname := "Ivan"
	role := 1
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	require.NoError(t, err)

	newNickname := "Paul"
	newRole := 0
	NewPasswordHash, err := bcrypt.GenerateFromPassword([]byte("87654321"), 10)
	require.NoError(t, err)

	user := domains.User{
		ID:           userID,
		Nickname:     newNickname,
		Role:         newRole,
		PasswordHash: NewPasswordHash,
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	u := &UserRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (id, nickname, role, password_hash) VALUES ($1, $2, $3, $4);
`, userID, nickname, role, passwordHash)
	require.NoError(t, err)

	err = u.UpdateUser(ctx, user)
	require.NoError(t, err)

	savedUser := &domains.User{}
	err = tx.QueryRow(ctx, `
	SELECT nickname, role, password_hash
	FROM users
	WHERE id=$1;
`, userID).Scan(&savedUser.Nickname, &savedUser.Role, &savedUser.PasswordHash)
	require.NoError(t, err)

	assert.Equal(t, newNickname, savedUser.Nickname)
	assert.Equal(t, newRole, savedUser.Role)
	assert.Equal(t, NewPasswordHash, savedUser.PasswordHash)
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	userID := uuid.New()
	nickname := "Ivan"
	role := 1
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	require.NoError(t, err)

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	u := &UserRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (id, nickname, role, password_hash) VALUES ($1, $2, $3, $4);
`, userID, nickname, role, passwordHash)
	require.NoError(t, err)

	err = u.DeleteUserByID(ctx, userID)
	require.NoError(t, err)

	savedUser := &domains.User{}
	err = tx.QueryRow(ctx, `
	SELECT nickname, role, password_hash
	FROM users
	WHERE id=$1;
`, userID).Scan(&savedUser.Nickname, &savedUser.Role, &savedUser.PasswordHash)
	require.Error(t, err)
}
