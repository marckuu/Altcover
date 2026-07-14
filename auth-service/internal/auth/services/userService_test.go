package services

import (
	"auth-service/core/domains"
	"auth-service/core/enums"
	messagingMocks "auth-service/core/shared/messaging/interfaces"
	repositoriesMocks "auth-service/internal/auth/repositories/interfaces"
	servicesMocks "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport/dto"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	userRepository := repositoriesMocks.NewMockUserRepository(t)
	tokenManager := repositoriesMocks.NewMockTokenManager(t)
	tokenService := servicesMocks.NewMockTokenService(t)
	producer := messagingMocks.NewMockProducer(t)
	ctx := context.Background()

	loginData := dto.LoginRequest{
		Nickname: "Ivan",
		Password: "12345678",
	}

	userRepository.
		On("GetUserByNickname", ctx, loginData.Nickname).
		Return(domains.User{}, errors.New("пользователь существует"))

	userRepository.
		On("AddUser", ctx, mock.MatchedBy(func(user domains.User) bool {
			if user.Nickname != loginData.Nickname {
				return false
			}
			if user.Role != enums.User {
				return false
			}
			return true
		})).
		Return(nil)

	u := &UserService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		tokenService:   tokenService,
		producer:       producer,
	}

	if err := u.Register(ctx, loginData); err != nil {
		t.Errorf("ошибка при регистрации: %v", err)
	}
}

func TestLogin(t *testing.T) {
	userRepository := repositoriesMocks.NewMockUserRepository(t)
	tokenManager := repositoriesMocks.NewMockTokenManager(t)
	tokenService := servicesMocks.NewMockTokenService(t)
	producer := messagingMocks.NewMockProducer(t)
	ctx := context.Background()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), 10)
	if err != nil {
		t.Fatalf("не удалось преобразовать пароль в хеш: %v", err)
	}

	loginData := dto.LoginRequest{
		Nickname: "Ivan",
		Password: "12345678",
	}

	tokenPair := repositoriesMocks.TokenPair{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	user := domains.User{
		ID:           uuid.New(),
		Nickname:     loginData.Nickname,
		Role:         1,
		PasswordHash: passwordHash,
		CreatedAt:    time.Time{},
	}

	userRepository.
		On("GetUserByNickname", ctx, loginData.Nickname).
		Return(user, nil)

	tokenManager.
		On("GenerateTokenPair", user.ID).
		Return(&tokenPair,
			nil)

	tokenService.
		On("AddRefreshToken", ctx, tokenPair.RefreshToken).
		Return(nil)

	u := &UserService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		tokenService:   tokenService,
		producer:       producer,
	}

	resTokenPair, err := u.Login(ctx, loginData)
	if err != nil {
		t.Errorf("ошибка в входа в аккаунт: %v", err)
	} else if resTokenPair.RefreshToken != tokenPair.RefreshToken && resTokenPair.AccessToken != tokenPair.AccessToken {
		t.Errorf("возвращены неккоректные токены: %v", err)
	}
}

func TestRefresh(t *testing.T) {
	userRepository := repositoriesMocks.NewMockUserRepository(t)
	tokenManager := repositoriesMocks.NewMockTokenManager(t)
	tokenService := servicesMocks.NewMockTokenService(t)
	producer := messagingMocks.NewMockProducer(t)
	ctx := context.Background()

	cookie := "refreshToken"

	userUuid := uuid.New()

	claims := &repositoriesMocks.CustomClaims{
		TokenType: "access",
		RegisteredClaims: repositoriesMocks.RegisteredClaims{
			Subject: userUuid.String(),
		},
	}

	tokenManager.
		On("Parse", cookie).
		Return(claims, nil).
		On("GenerateAccessToken", userUuid).
		Return("accessToken", nil)

	tokenService.
		On("IsTokenRevoked", ctx, cookie).
		Return(false, nil)

	u := &UserService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		tokenService:   tokenService,
		producer:       producer,
	}

	token, err := u.Refresh(ctx, cookie)
	if err != nil {
		t.Errorf("не удалось обновить токен: %v", err)
	}

	if token == "" {
		t.Error("возвращен пустой токен")
	}
}

func TestLogout(t *testing.T) {
	userRepository := repositoriesMocks.NewMockUserRepository(t)
	tokenManager := repositoriesMocks.NewMockTokenManager(t)
	tokenService := servicesMocks.NewMockTokenService(t)
	producer := messagingMocks.NewMockProducer(t)
	ctx := context.Background()

	cookie := "refreshToken"

	tokenManager.
		On("Parse", cookie).
		Return(nil, nil)

	tokenService.
		On("DeleteRefreshToken", ctx, cookie).
		Return(nil)

	u := &UserService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		tokenService:   tokenService,
		producer:       producer,
	}

	if err := u.Logout(ctx, cookie); err != nil {
		t.Errorf("ошибка выхода из аккаунта: %v", err)
	}
}
