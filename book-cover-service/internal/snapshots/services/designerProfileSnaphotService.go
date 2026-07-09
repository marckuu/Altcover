package services

import (
	"book-cover-service/internal/snapshots/repositories/interfaces"
	"book-cover-service/internal/snapshots/transport/dto"
	"context"

	"github.com/google/uuid"
)

type DesignerProfileSnapshotService struct {
	designerProfileSnapshotRepository interfaces.DesignerProfileSnapshotRepository
}

func NewDesignerProfileSnapshotService(designerProfileSnapshotRepository interfaces.DesignerProfileSnapshotRepository) *DesignerProfileSnapshotService {
	return &DesignerProfileSnapshotService{
		designerProfileSnapshotRepository: designerProfileSnapshotRepository,
	}
}

func (d *DesignerProfileSnapshotService) AddDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error {
	if err := d.designerProfileSnapshotRepository.AddDesignerProfileSnapshot(ctx, profile); err != nil {
		return err
	}
	return nil
}

func (d *DesignerProfileSnapshotService) UpdateDesignerProfileSnapshot(ctx context.Context, profile dto.DesignerProfileSnapshot) error {
	if err := d.designerProfileSnapshotRepository.UpdateDesignerProfileSnapshot(ctx, profile); err != nil {
		return err
	}
	return nil
}

func (d *DesignerProfileSnapshotService) DeleteDesignerProfileSnapshot(ctx context.Context, profileID uuid.UUID) error {
	if err := d.designerProfileSnapshotRepository.DeleteDesignerProfileSnapshot(ctx, profileID); err != nil {
		return err
	}
	return nil
}

func (d *DesignerProfileSnapshotService) GetDesignerProfileSnapshotByUserID(ctx context.Context, userID uuid.UUID) (dto.DesignerProfileSnapshot, error) {
	profileSnapshot, err := d.designerProfileSnapshotRepository.GetDesignerProfileSnapshotByUserID(ctx, userID)
	if err != nil {
		return dto.DesignerProfileSnapshot{}, err
	}
	return profileSnapshot, nil
}
