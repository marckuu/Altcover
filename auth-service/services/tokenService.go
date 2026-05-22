package services

import (
	"auth-service/repositories"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	tokenRepository repositories.TokenRepository
}

func (t *TokenService) AddRefreshToken(ctx context.Context, token string) error {
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 1)
	if err != nil {
		return fmt.Errorf("не удалось получить хэш refresh токена, %w", err)
	}

	if err = t.tokenRepository.AddRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("не удалось сохранить хэш refresh токена, %w", err)
	}

	return nil
}

func (t *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 1)
	if err != nil {
		return fmt.Errorf("не удалось получить хэш refresh токена, %w", err)
	}

	if err = t.tokenRepository.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("не удалось удалить хэш refresh токена, %w", err)
	}

	return nil
}

func (t *TokenService) IsTokenRevoked(ctx context.Context, tokenRaw string) (bool, error) {
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(tokenRaw), 1)
	if err != nil {
		return false, fmt.Errorf("не удалось получить хэш refresh токена, %w", err)
	}

	_, err = t.tokenRepository.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repositories.ErrTokenNotFound) {
			return true, nil
		}
		return false, fmt.Errorf("ошибка получения хэша refresh токена из базы, %w", err)
	}

	return false, nil
}
