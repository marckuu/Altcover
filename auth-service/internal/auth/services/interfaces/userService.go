package interfaces

import (
	globInterfaces "auth-service/internal/auth/repositories/interfaces"
	"auth-service/internal/auth/transport/dto"
	"context"
)

type UserService interface {
	Register(ctx context.Context, loginRequest dto.LoginRequest) error
	Login(ctx context.Context, loginRequest dto.LoginRequest) (*globInterfaces.TokenPair, error)
	Refresh(ctx context.Context, cookieValue string) (string, error)
	Logout(ctx context.Context, cookieValue string) error
}
