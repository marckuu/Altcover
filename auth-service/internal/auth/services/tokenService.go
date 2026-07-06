package services

import (
	globalErrors "auth-service/core/errors"
	"auth-service/core/shared/messaging"
	"auth-service/core/tools"
	"auth-service/internal/auth/repositories/interfaces"
	"context"
	"errors"
	"fmt"
)

type TokenService struct {
	tokenRepository interfaces.TokenRepository
	producer        *messaging.Producer
}

func NewTokenService(tokenRepository interfaces.TokenRepository, producer *messaging.Producer) *TokenService {
	return &TokenService{
		tokenRepository: tokenRepository,
		producer:        producer,
	}
}

func (t *TokenService) AddRefreshToken(ctx context.Context, token string) error {
	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		return err
	}

	if err = t.tokenRepository.AddRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("не удалось сохранить хэш refresh токена: %w", err)
	}

	return nil
}

func (t *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		return err
	}

	if err = t.tokenRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("не удалось удалить хэш refresh токена, %w", err)
	}

	return nil
}

func (t *TokenService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		return false, err
	}

	_, err = t.tokenRepository.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, globalErrors.ErrTokenNotFound) {
			return true, nil
		}
		return false, fmt.Errorf("ошибка получения хэша refresh токена из базы, %w", err)
	}

	return false, nil
}
