package interfaces

import (
	"auth-service/core/domains"
	globInterfaces "auth-service/internal/auth/repositories/interfaces"
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	AddUser(ctx context.Context, user domains.User) error
	GetUser() domains.User
	GetUserIDFromTokenClaims(claims *globInterfaces.CustomClaims) (uuid.UUID, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (domains.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (domains.User, error)
}
