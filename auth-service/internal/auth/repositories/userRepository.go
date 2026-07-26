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

var errUserNotFound = errors.New("пользователь не найден")

type UserRepository struct {
	connection db.Database
}

func NewUserRepository(conn db.Database) *UserRepository {
	return &UserRepository{
		connection: conn,
	}
}

func (u *UserRepository) AddUser(ctx context.Context, user domains.User) error {
	query := `
	INSERT INTO users (nickname, role, password_hash)
	VALUES ($1, $2, $3);
`
	if _, err := u.connection.Exec(ctx, query, user.Nickname, user.Role, user.PasswordHash); err != nil {
		return fmt.Errorf("user repo -> add: %w", err)
	}

	return nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (domains.User, error) {
	query := `
	SELECT id, nickname, role, created_at, password_hash 
	FROM users
	WHERE id = $1;
`
	resultRow := u.connection.QueryRow(ctx, query, userID)

	user := domains.User{}

	if err := resultRow.Scan(&user.ID, &user.Nickname, &user.Role, &user.CreatedAt, &user.PasswordHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.User{}, errUserNotFound
		}
		return domains.User{}, fmt.Errorf("user repo -> get by id: %w", err)
	}

	return user, nil
}

func (u *UserRepository) GetUsers(ctx context.Context, usersIDs []int64, offset int, limit int) ([]domains.User, error) {
	query := `
	SELECT id, nickname, role, created_at, password_hash
	FROM users
	WHERE id = ANY($1)
	LIMIT $2
	OFFSET $3;
`
	resultRows, err := u.connection.Query(ctx, query, usersIDs, limit, offset)
	if err != nil {
		return []domains.User{}, fmt.Errorf("user repo -> get all -> query: %w", err)
	}

	defer resultRows.Close()

	var users []domains.User

	for resultRows.Next() {
		user := domains.User{}
		if err = resultRows.Scan(&user.ID, &user.Nickname, &user.Role, &user.CreatedAt, &user.PasswordHash); err != nil {
			return []domains.User{}, fmt.Errorf("user repo -> get all -> parsing: %w", err)
		}
		users = append(users, user)
	}

	if err = resultRows.Err(); err != nil {
		return []domains.User{}, fmt.Errorf("user repo -> get all -> query thread: %w", err)
	}

	return users, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user domains.User) error {
	query := `
	UPDATE users
	SET nickname = $1,
	    role = $2,
	    password_hash = $3
	WHERE id = $4;
`
	tag, err := u.connection.Exec(ctx, query, user.Nickname, user.Role, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("user repo -> update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errUserNotFound
	}

	return nil
}

func (u *UserRepository) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	query := `
	DELETE FROM users
	WHERE id = $1;
`
	tag, err := u.connection.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("user repo -> delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errUserNotFound
	}

	return nil
}

func (u *UserRepository) GetUserByNickname(ctx context.Context, nickname string) (domains.User, error) {
	query := `
	SELECT id, nickname, role, password_hash, created_at
	FROM users
	WHERE nickname = $1;
`
	resultRow := u.connection.QueryRow(ctx, query, nickname)

	var user domains.User

	if err := resultRow.Scan(&user.ID, &user.Nickname, &user.Role, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.User{}, errUserNotFound
		}
		return domains.User{}, err
	}
	return user, nil
}
