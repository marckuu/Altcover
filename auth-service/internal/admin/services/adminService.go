package services

import (
	"auth-service/core/enums"
	"auth-service/internal/admin/repositories/interfaces"
	"context"

	"github.com/google/uuid"
)

type AdminService struct {
	adminRepository interfaces.AdminRepository
}

func NewAdminService(adminRepository interfaces.AdminRepository) AdminService {
	return AdminService{
		adminRepository: adminRepository,
	}
}

func (a AdminService) ChangeRole(ctx context.Context, userID uuid.UUID, userRole enums.Role) error {
	if err := a.adminRepository.UpdateRole(ctx, userID, userRole); err != nil {
		return err
	}
	return nil
}

func (a AdminService) GetRole(ctx context.Context, userID uuid.UUID) (enums.Role, error) {
	role, err := a.adminRepository.GetRole(ctx, userID)
	if err != nil {
		return enums.User, err
	}
	return role, nil
}
