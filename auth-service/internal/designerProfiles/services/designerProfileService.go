package services

import (
	"auth-service/core/domains"
	"auth-service/core/shared"
	"auth-service/core/shared/messaging"
	"auth-service/internal/designerProfiles/repositories/interfaces"
	"context"
	"encoding/json"
	"fmt"

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

func (d *DesignerProfileService) CreateDesignerProfileToUser(ctx context.Context, userID uuid.UUID, designerProfile domains.DesignerProfile) (domains.DesignerProfile, error) {
	_, err := d.designerProfileRepository.GetDesignerProfileByUserID(ctx, userID)
	if err == nil {
		return domains.DesignerProfile{}, fmt.Errorf("профиль дизайнера уже существует: %w", err)
	}

	savedProfile, err := d.designerProfileRepository.AddDesignerProfile(ctx, designerProfile)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	payload, err := json.Marshal(designerProfile)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	event := shared.NewKafkaEvent(shared.DesignerProfileCreated, payload)

	message, err := json.Marshal(event)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	if err = d.producer.Produce(message); err != nil {
		return domains.DesignerProfile{}, err
	}

	return savedProfile, nil
}

func (d *DesignerProfileService) UpdateDesignerProfileToUser(ctx context.Context, userID uuid.UUID, newDesignerProfile domains.DesignerProfile) (domains.DesignerProfile, error) {
	oldDesignerProfile, err := d.designerProfileRepository.GetDesignerProfileByUserID(ctx, userID)
	if err != nil {
		return domains.DesignerProfile{}, fmt.Errorf("ошибка получения профиля дизайнера: %w", err)
	}

	newDesignerProfile.ID = oldDesignerProfile.ID

	savedProfile, err := d.designerProfileRepository.UpdateDesignerProfile(ctx, newDesignerProfile)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	payload, err := json.Marshal(newDesignerProfile)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	event := shared.NewKafkaEvent(shared.DesignerProfileUpdated, payload)

	message, err := json.Marshal(event)
	if err != nil {
		return domains.DesignerProfile{}, err
	}

	if err = d.producer.Produce(message); err != nil {
		return domains.DesignerProfile{}, err
	}

	return savedProfile, nil
}

func (d *DesignerProfileService) DeleteDesignerProfileToUser(ctx context.Context, userID uuid.UUID) error {
	designerProfile, err := d.designerProfileRepository.GetDesignerProfileByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("не удалось получить профиль дизайнера: %w", err)
	}

	if err = d.designerProfileRepository.DeleteDesignerProfile(ctx, designerProfile.ID); err != nil {
		return fmt.Errorf("не удалось удалить профиль: %w", err)
	}

	payload, err := json.Marshal(domains.DesignerProfile{
		ID: designerProfile.ID,
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
