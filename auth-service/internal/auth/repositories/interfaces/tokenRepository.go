package interfaces

import (
	"auth-service/core/domains"
	"context"
)

type TokenRepository interface {
	AddRefreshToken(ctx context.Context, tokenHash []byte) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash []byte) (domains.Token, error)
	DeleteRefreshToken(ctx context.Context, tokenHash []byte) error
}
