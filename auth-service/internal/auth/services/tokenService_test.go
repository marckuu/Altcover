package services

import (
	"auth-service/core/domains"
	globalErrors "auth-service/core/errors"
	messagingInterfaces "auth-service/core/shared/messaging/interfaces"
	"auth-service/core/tools"
	"auth-service/internal/auth/repositories/interfaces"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestTable struct {
	testName     string
	TokenService *TokenService
	Expected     bool
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

	err = s.AddRefreshToken(ctx, token)
	assert.NoError(t, err)
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

	err = s.DeleteRefreshToken(ctx, token)
	assert.NoError(t, err)
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
			"true",
			&TokenService{
				tokenRepository: tokenRepositoryTrue,
				producer:        producer,
			},
			true,
		},
		{
			"false",
			&TokenService{
				tokenRepository: tokenRepositoryFalse,
				producer:        producer,
			},
			false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			revoked, err := tt.TokenService.IsTokenRevoked(ctx, token)

			require.NoError(t, err)
			assert.Equal(t, tt.Expected, revoked)
		})
	}
}
