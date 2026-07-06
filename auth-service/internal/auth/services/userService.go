package services

import (
	"auth-service/core/domains"
	"auth-service/core/shared/messaging"
	"auth-service/internal/auth/repositories/interfaces"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	userRepository interfaces.UserRepository
	producer       *messaging.Producer
}

func NewUserService(repository interfaces.UserRepository, producer *messaging.Producer) *UserService {
	return &UserService{
		userRepository: repository,
		producer:       producer,
	}
}

func (u *UserService) AddUser(ctx context.Context, user domains.User) error {
	if err := u.userRepository.AddUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUser() domains.User {
	// 1. Validate

	// 2. repo.Get()

	return domains.User{}
}

func (u *UserService) UpdateUser(user domains.User) {
	// 1. Validate

	// 2. repo.Update()
}

func (u *UserService) DeleteUser(userID pgtype.UUID) {
	// 1. Validate

	// 2. repo.Update()
}

func (u *UserService) GetUserIDFromTokenClaims(claims *interfaces.CustomClaims) (uuid.UUID, error) {
	userID := claims.RegisteredClaims.Subject
	userIDConverted, err := uuid.Parse(userID)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userIDConverted, nil
}

func (u *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (domains.User, error) {
	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return domains.User{}, err
	}
	return user, nil
}

func (u *UserService) GetUserByNickname(ctx context.Context, nickname string) (domains.User, error) {
	user, err := u.userRepository.GetUserByNickname(ctx, nickname)
	if err != nil {
		return domains.User{}, err
	}
	return user, nil
}
