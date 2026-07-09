package interfaces

import (
	"auth-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domains.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (domains.User, error)
	GetUsers(ctx context.Context, usersIDs []int64, offset int, limit int) ([]domains.User, error)
	UpdateUser(ctx context.Context, user domains.User) error
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
	GetUserByNickname(ctx context.Context, nickname string) (domains.User, error)
}
