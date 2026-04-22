package services

import (
	"Altcover/auth-service/domains"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	userRepository UserRepository
}

type UserRepository interface {
}

func (us *UserService) CreateUser(user domains.User) {
	// 1. Validate

	// 2. repo.Create()
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
