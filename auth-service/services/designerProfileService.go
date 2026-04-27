package services

import (
	"auth-service/domains"
	"auth-service/repositories"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileService struct {
	designerProfileRepository repositories.DesignerProfileRepository
}

func (ds *DesignerProfileService) UpdateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error {

	// 1. Validate

	// 2. repo.Update()
	if err := ds.designerProfileRepository.UpdateDesignerProfile(ctx, designerProfile); err != nil {
		return err
	}

	return nil
}

func (ds *DesignerProfileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error) {
	designerProfile, err := ds.designerProfileRepository.GetDesignerProfileByUserID(ctx, userID)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	return designerProfile, nil
}
