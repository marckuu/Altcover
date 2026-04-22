package services

import (
	"Altcover/auth-service/domains"

	"github.com/jackc/pgx/v5/pgtype"
)

type DesignerProfileService struct {
	designerProfileRepository DesignerProfileRepository
}

type DesignerProfileRepository interface {
}

func (ds *DesignerProfileService) CreateDesignerProfile(designerProfile domains.DesignerProfile) {
	// 1. Validate

	// 2. repo.Create()
}

func (ds *DesignerProfileService) GetDesignerProfile(designerProfileID pgtype.UUID) domains.DesignerProfile {
	// 1. Validate

	// 2. repo.Get()

	return domains.DesignerProfile{}
}

func (ds *DesignerProfileService) UpdateDesignerProfile(designerProfile domains.DesignerProfile) {
	// 1. Validate

	// 2. repo.Update()
}
