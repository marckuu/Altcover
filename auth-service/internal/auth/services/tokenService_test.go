package services

import (
	"auth-service/core/domains"
	globalErrors "auth-service/core/errors"
	messagingInterfaces "auth-service/core/shared/messaging/interfaces"
	"auth-service/core/tools"
	"auth-service/internal/auth/repositories/interfaces"
	"context"
	"testing"
)

type TestTable struct {
	TokenService *TokenService
	Result       bool
}

func TestAddRefreshToken(t *testing.T) {
	tokenRepository := interfaces.NewMockTokenRepository(t)
	producer := messagingInterfaces.NewMockProducer(t)
	ctx := context.Background()
	token := "token"

	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		t.Fatalf("ошибка получения хеша токена:, %v", err)
	}

	tokenRepository.
		On("AddRefreshToken", ctx, tokenHash).
		Return(nil)

	s := &TokenService{
		tokenRepository: tokenRepository,
		producer:        producer,
	}

	if err = s.AddRefreshToken(ctx, token); err != nil {
		t.Errorf("ошибка добавления refresh токена:, %d", err)
	}
}

func TestDeleteRefreshToken(t *testing.T) {
	tokenRepository := interfaces.NewMockTokenRepository(t)
	producer := messagingInterfaces.NewMockProducer(t)
	ctx := context.Background()
	token := "token"

	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		t.Fatalf("ошибка получения хеша токена:, %v", err)
	}

	tokenRepository.
		On("DeleteRefreshToken", ctx, tokenHash).
		Return(nil)

	s := &TokenService{
		tokenRepository: tokenRepository,
		producer:        producer,
	}

	if err = s.DeleteRefreshToken(ctx, token); err != nil {
		t.Errorf("ошибка удаления refresh токена: %v", err)
	}

}

func TestIsTokenRevoked(t *testing.T) {
	producer := messagingInterfaces.NewMockProducer(t)
	ctx := context.Background()
	token := "token"

	tokenHash, err := tools.GetTokenHash(token)
	if err != nil {
		t.Fatalf("ошибка получения хеша токена:, %v", err)
	}

	tokenRepositoryTrue := interfaces.NewMockTokenRepository(t)
	tokenRepositoryFalse := interfaces.NewMockTokenRepository(t)

	tokenRepositoryTrue.
		On("GetRefreshTokenByHash", ctx, tokenHash).
		Return(domains.Token{}, globalErrors.ErrTokenNotFound)

	tokenRepositoryFalse.
		On("GetRefreshTokenByHash", ctx, tokenHash).
		Return(domains.Token{}, nil)

	testTable := []TestTable{
		{
			&TokenService{
				tokenRepository: tokenRepositoryTrue,
				producer:        producer,
			},
			true,
		},
		{
			&TokenService{
				tokenRepository: tokenRepositoryFalse,
				producer:        producer,
			},
			false,
		},
	}

	for _, tt := range testTable {
		revoked, err := tt.TokenService.IsTokenRevoked(ctx, token)
		if err != nil {
			t.Errorf("ошибка определения отозван ли токен: %v", err)
		}

		if revoked != tt.Result {
			t.Errorf("неккоректный результат, revoked: %v", err)
		}
	}
}
