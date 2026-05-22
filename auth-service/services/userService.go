package services

import (
	"auth-service/domains"
	"auth-service/repositories"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	userRepository repositories.UserRepository
}

type UserRepository interface {
}

func (us *UserService) AddUser(ctx context.Context, user domains.User) error {
	if err := us.userRepository.AddUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUser() domains.User {
	// 1. Validate

	// 2. repo.Get()

	return domains.User{}
}

func (us *UserService) UpdateUser(user domains.User) {
	// 1. Validate

	// 2. repo.Update()
}

func (us *UserService) DeleteUser(userID pgtype.UUID) {
	// 1. Validate

	// 2. repo.Update()
}

func (u *UserService) GetUserIDFromTokenClaims(claims *repositories.CustomClaims) (uuid.UUID, error) {
	userID := claims.Subject
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
