package repositories

import (
	"auth-service/core/domains"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testTimeout = 10 * time.Second

func TestDesignerProfileRepository_AddDesignerProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	profile := domains.DesignerProfile{
		Nickname:  "Ivan",
		AvatarKey: "12345678",
	}
	user := domains.User{
		Nickname:     "Ivan",
		Role:         1,
		PasswordHash: []byte{1, 0},
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	d := &DesignerProfileRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (nickname, role, password_hash)
	VALUES ($1, $2, $3);
`, user.Nickname, user.Role, user.PasswordHash)
	require.NoError(t, err)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id
	FROM users
	WHERE nickname=$1;
`, user.Nickname).Scan(&userID)
	require.NoError(t, err)

	profile.UserID = userID

	err = d.AddDesignerProfile(ctx, profile)
	require.NoError(t, err)

	receivedProfile := domains.DesignerProfile{}
	err = tx.QueryRow(ctx, `
	SELECT nickname, avatar_key, user_id 
	FROM designer_profile
	WHERE nickname=$1;
`, &profile.Nickname).Scan(&receivedProfile.Nickname, &receivedProfile.AvatarKey, &receivedProfile.UserID)
	require.NoError(t, err)

	assert.Equal(t, profile.Nickname, receivedProfile.Nickname)
	assert.Equal(t, profile.AvatarKey, receivedProfile.AvatarKey)
	assert.Equal(t, profile.UserID, receivedProfile.UserID)
}

func TestDesignerProfileRepository_GetDesignerProfileByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	profile := domains.DesignerProfile{
		Nickname:  "Ivan",
		AvatarKey: "12345678",
	}
	user := domains.User{
		Nickname:     "Ivan",
		Role:         1,
		PasswordHash: []byte{1, 0},
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	d := &DesignerProfileRepository{
		connection: tx,
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (nickname, role, password_hash)
	VALUES ($1, $2, $3);
`, &user.Nickname, &user.Role, &user.PasswordHash)
	require.NoError(t, err)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id
	FROM users
	WHERE nickname=$1;
`, user.Nickname).Scan(&userID)
	require.NoError(t, err)

	_, err = tx.Exec(ctx, `
	INSERT INTO designer_profile (nickname, avatar_key, user_id)
	VALUES ($1, $2, $3);
`, profile.Nickname, profile.AvatarKey, userID)
	require.NoError(t, err)

	receivedProfile, err := d.GetDesignerProfileByUserID(ctx, userID)
	require.NoError(t, err)

	assert.Equal(t, profile.Nickname, receivedProfile.Nickname)
	assert.Equal(t, profile.AvatarKey, receivedProfile.AvatarKey)
	assert.Equal(t, userID, receivedProfile.UserID)
}

func TestDesignerProfileRepository_UpdateDesignerProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	profile := domains.DesignerProfile{
		Nickname:  "Ivan",
		AvatarKey: "12345678",
	}
	newProfile := domains.DesignerProfile{
		Nickname:  "Ron",
		AvatarKey: "87654321",
	}
	user := domains.User{
		Nickname:     "Ivan",
		Role:         1,
		PasswordHash: []byte{1, 0},
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	d := &DesignerProfileRepository{connection: tx}

	_, err = tx.Exec(ctx, `
	INSERT INTO users (nickname, role, password_hash)
	VALUES ($1, $2, $3);
`, &user.Nickname, &user.Role, &user.PasswordHash)
	require.NoError(t, err)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id
	FROM users
	WHERE nickname=$1;
`, user.Nickname).Scan(&userID)
	require.NoError(t, err)

	_, err = tx.Exec(ctx, `
	INSERT INTO designer_profile (nickname, avatar_key, user_id)
	VALUES ($1, $2, $3);
`, profile.Nickname, profile.AvatarKey, userID)
	require.NoError(t, err)

	var profileID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id 
	FROM designer_profile
	WHERE nickname=$1;
`, &profile.Nickname).Scan(&profileID)
	require.NoError(t, err)

	newProfile.ID = profileID

	err = d.UpdateDesignerProfile(ctx, newProfile)
	require.NoError(t, err)

	receivedProfile := domains.DesignerProfile{}
	err = tx.QueryRow(ctx, `
	SELECT nickname, avatar_key, user_id 
	FROM designer_profile
	WHERE nickname=$1;
`, &newProfile.Nickname).Scan(&receivedProfile.Nickname, &receivedProfile.AvatarKey, &receivedProfile.UserID)
	require.NoError(t, err)

	assert.Equal(t, newProfile.Nickname, receivedProfile.Nickname)
	assert.Equal(t, newProfile.AvatarKey, receivedProfile.AvatarKey)
	assert.Equal(t, userID, receivedProfile.UserID)
}

func TestDesignerProfileRepository_DeleteDesignerProfile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	profile := domains.DesignerProfile{
		Nickname:  "Ivan",
		AvatarKey: "12345678",
	}
	user := domains.User{
		Nickname:     "Ivan",
		Role:         1,
		PasswordHash: []byte{1, 0},
	}

	tx, err := conn.Begin(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		cancel()
	})

	_, err = tx.Exec(ctx, `
	INSERT INTO users (nickname, role, password_hash)
	VALUES ($1, $2, $3);
`, &user.Nickname, &user.Role, &user.PasswordHash)
	require.NoError(t, err)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id
	FROM users
	WHERE nickname=$1;
`, user.Nickname).Scan(&userID)
	require.NoError(t, err)

	_, err = tx.Exec(ctx, `
	INSERT INTO designer_profile (nickname, avatar_key, user_id)
	VALUES ($1, $2, $3);
`, profile.Nickname, profile.AvatarKey, userID)
	require.NoError(t, err)

	var profileID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id 
	FROM designer_profile
	WHERE nickname=$1;
`, &profile.Nickname).Scan(&profileID)
	require.NoError(t, err)

	d := &DesignerProfileRepository{connection: tx}

	err = d.DeleteDesignerProfile(ctx, profileID)
	require.NoError(t, err)

	var receivedProfileID uuid.UUID
	err = tx.QueryRow(ctx, `
	SELECT id 
	FROM designer_profile
	WHERE nickname=$1;
`, &profile.Nickname).Scan(&receivedProfileID)
	assert.Equal(t, err, pgx.ErrNoRows)

}
