package services

import (
	"auth-service/core/domains"
	"auth-service/core/enums"
	"auth-service/core/shared/messaging"
	authRepositoryInterfaces "auth-service/internal/auth/repositories/interfaces"
	authServicesInterfaces "auth-service/internal/auth/services/interfaces"
	"auth-service/internal/auth/transport/dto"
	"context"

	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository authRepositoryInterfaces.UserRepository
	tokenManager   authRepositoryInterfaces.TokenManager
	tokenService   authServicesInterfaces.TokenService
	producer       *messaging.Producer
}

func NewUserService(repository authRepositoryInterfaces.UserRepository,
	producer *messaging.Producer,
	tokenManager authRepositoryInterfaces.TokenManager,
	tokenService authServicesInterfaces.TokenService) *UserService {
	return &UserService{
		userRepository: repository,
		tokenManager:   tokenManager,
		tokenService:   tokenService,
		producer:       producer,
	}
}

func (u *UserService) Register(ctx context.Context, loginRequest dto.LoginRequest) error {
	_, err := u.userRepository.GetUserByNickname(ctx, loginRequest.Nickname)
	if err == nil {
		return fmt.Errorf("пользователь уже существует: %w", err)
	}

	PasswordHash, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), 10)
	if err != nil {
		return fmt.Errorf("ошибка получения хэша пароля: %w", err)
	}

	var user domains.User

	user.Role = enums.User
	user.Nickname = loginRequest.Nickname
	user.PasswordHash = PasswordHash

	if err = u.userRepository.AddUser(ctx, user); err != nil {
		return fmt.Errorf("ошибка добавления нового пользователя: %w", err)
	}

	return nil

}

func (u *UserService) Login(ctx context.Context, loginRequest dto.LoginRequest) (*authRepositoryInterfaces.TokenPair, error) {
	user, err := u.userRepository.GetUserByNickname(ctx, loginRequest.Nickname)
	if err != nil {
		return &authRepositoryInterfaces.TokenPair{}, fmt.Errorf("пользователь не найден: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginRequest.Password)); err != nil {
		return &authRepositoryInterfaces.TokenPair{}, fmt.Errorf("неверный пароль: %w", err)
	}

	tokenPair, err := u.tokenManager.GenerateTokenPair(user.ID)
	if err != nil {
		return &authRepositoryInterfaces.TokenPair{}, fmt.Errorf("не удалось сгенерировать jwtCovers токены: %w", err)
	}

	if err = u.tokenService.AddRefreshToken(ctx, tokenPair.RefreshToken); err != nil {
		return &authRepositoryInterfaces.TokenPair{}, fmt.Errorf("не удалось сохранить refresh токен: %w", err)
	}

	return tokenPair, nil

}

func (u *UserService) Refresh(ctx context.Context, cookieValue string) (string, error) {
	claims, err := u.tokenManager.Parse(cookieValue)
	if err != nil {
		return "", fmt.Errorf("неккоректный refresh токен: %w", err)
	}

	isRevoked, err := u.tokenService.IsTokenRevoked(ctx, cookieValue)
	if err != nil {
		return "", fmt.Errorf("не удалось проверить токен: %w", err)
	}

	if isRevoked {
		return "", fmt.Errorf("переданный refresh токен отозван: %w", err)
	}

	userID, err := uuid.Parse(claims.RegisteredClaims.Subject)
	if err != nil {
		return "", fmt.Errorf("ошибка получения id пользователя из refresh токена: %w", err)
	}

	accessToken, err := u.tokenManager.GenerateAccessToken(userID)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать access токен: %w", err)
	}

	return accessToken, nil
}

func (u *UserService) Logout(ctx context.Context, cookieValue string) error {
	_, err := u.tokenManager.Parse(cookieValue)
	if err != nil {
		return fmt.Errorf("неккоректный refresh токен: %w", err)
	}

	if err = u.tokenService.DeleteRefreshToken(ctx, cookieValue); err != nil {
		return fmt.Errorf("ошибка удаления refresh токена из базы: %w", err)
	}

	return nil
}
