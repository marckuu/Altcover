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

func (d *DesignerProfileService) CreateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error {
	if err := d.designerProfileRepository.AddDesignerProfile(ctx, designerProfile); err != nil {
		return err
	}

	return nil
}

func (d *DesignerProfileService) UpdateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error {

	// 1. Validate

	// 2. repo.Update()
	if err := d.designerProfileRepository.UpdateDesignerProfile(ctx, designerProfile); err != nil {
		return err
	}

	return nil
}

func (d *DesignerProfileService) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domains.DesignerProfile, error) {
	designerProfile, err := d.designerProfileRepository.GetDesignerProfileByUserID(ctx, userID)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	return designerProfile, nil
}

func (d *DesignerProfileService) DeleteDesignerProfile(ctx context.Context, profileID uuid.UUID) error {
	if err := d.designerProfileRepository.DeleteDesignerProfile(ctx, profileID); err != nil {
		return err
	}
	return nil
}
