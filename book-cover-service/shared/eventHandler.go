package shared

import (
	"book-cover-service/db/repositories"
	"book-cover-service/dto"
	"context"
	"encoding/json"
	"fmt"
)

type EventHandlers struct {
	designerProfileSnapshotRepository repositories.DesignerProfileSnapshotRepository
	ctx                               context.Context
}

func NewEventHandlers(ctx context.Context, repository repositories.DesignerProfileSnapshotRepository) EventHandlers {
	return EventHandlers{
		designerProfileSnapshotRepository: repository,
		ctx:                               ctx,
	}
}

func (e *EventHandlers) HandleKafkaEvent(message []byte) {
	var kafkaEvent KafkaEvent
	if err := json.Unmarshal(message, &kafkaEvent); err != nil {
		fmt.Printf("не удалось распрасить сообщение из kafka: %v", err)
		return
	}

	switch kafkaEvent.EventType {
	case designerProfileCreated:
		profile, err := e.DecodeDomain(kafkaEvent.Payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		if err = e.designerProfileSnapshotRepository.AddDesignerProfileSnapshot(e.ctx, profile); err != nil {
			fmt.Printf(err.Error())
		}
	case designerProfileUpdated:
		profile, err := e.DecodeDomain(kafkaEvent.Payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		if err = e.designerProfileSnapshotRepository.UpdateDesignerProfileSnapshot(e.ctx, profile); err != nil {
			fmt.Printf(err.Error())
		}
	case designerProfileDeleted:
		profile, err := e.DecodeDomain(kafkaEvent.Payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		if err = e.designerProfileSnapshotRepository.DeleteDesignerProfileSnapshot(e.ctx, profile.ID); err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func (e *EventHandlers) DecodeDomain(message json.RawMessage) (dto.DesignerProfileSnapshot, error) {
	var designerProfileSnapshot dto.DesignerProfileSnapshot

	if err := json.Unmarshal(message, &designerProfileSnapshot); err != nil {
		return dto.DesignerProfileSnapshot{}, fmt.Errorf("не удалось распрасить сущность внутри сообщения из kafka: %w", err)
	}

	return designerProfileSnapshot, nil
}
