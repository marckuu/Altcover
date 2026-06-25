package repositories

import (
	"book-cover-service/dto"

	"github.com/google/uuid"
)

type DesignerProfileSnapshotRepository struct {
}

func NewDesignerProfileRepository() DesignerProfileSnapshotRepository {
	return DesignerProfileSnapshotRepository{}
}

func (d *DesignerProfileSnapshotRepository) CreateDesignerProfileSnapshot(profile dto.DesignerProfileSnapshot) {

}

func (d *DesignerProfileSnapshotRepository) UpdateDesignerProfileSnapshot(profile dto.DesignerProfileSnapshot) {

}

func (d *DesignerProfileSnapshotRepository) DeleteDesignerProfileSnapshot(profileID uuid.UUID) {

}
