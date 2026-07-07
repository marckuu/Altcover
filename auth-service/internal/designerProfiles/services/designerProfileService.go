package services

import (
	"auth-service/core/domains"
	"auth-service/core/shared"
	"auth-service/core/shared/messaging"
	"auth-service/internal/designerProfiles/repositories/interfaces"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type DesignerProfileService struct {
	designerProfileRepository interfaces.DesignerProfileRepository
	producer                  *messaging.Producer
}

func NewDesignerProfileService(repository interfaces.DesignerProfileRepository, producer *messaging.Producer) *DesignerProfileService {
	return &DesignerProfileService{
		designerProfileRepository: repository,
		producer:                  producer,
	}
}

func (d *DesignerProfileService) CreateDesignerProfile(ctx context.Context, designerProfile domains.DesignerProfile) error {
	if err := d.designerProfileRepository.AddDesignerProfile(ctx, designerProfile); err != nil {
		return err
	}

	payload, err := json.Marshal(designerProfile)
	if err != nil {
		return err
	}

	event := shared.NewKafkaEvent(shared.DesignerProfileCreated, payload)

	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err = d.producer.Produce(message); err != nil {
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

	payload, err := json.Marshal(designerProfile)
	if err != nil {
		return err
	}

	event := shared.NewKafkaEvent(shared.DesignerProfileUpdated, payload)

	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err = d.producer.Produce(message); err != nil {
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

	payload, err := json.Marshal(domains.DesignerProfile{
		ID: profileID,
	})
	if err != nil {
		return err
	}

	event := shared.NewKafkaEvent(shared.DesignerProfileDeleted, payload)

	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err = d.producer.Produce(message); err != nil {
		return err
	}

	return nil
}
