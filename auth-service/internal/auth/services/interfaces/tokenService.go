package interfaces

import "context"

type TokenService interface {
	AddRefreshToken(ctx context.Context, token string) error
	DeleteRefreshToken(ctx context.Context, token string) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}
