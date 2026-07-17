package interfaces

import (
	"auth-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	AddUser(ctx context.Context, user domains.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (domains.User, error)
	UpdateUser(ctx context.Context, user domains.User) error
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
	GetUserByNickname(ctx context.Context, nickname string) (domains.User, error)
}
