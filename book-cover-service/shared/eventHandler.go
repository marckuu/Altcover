package shared

import (
	"book-cover-service/db/repositories"
	"book-cover-service/dto"
	"encoding/json"
	"fmt"
)

type EventHandlers struct {
	designerProfileSnapshotRepository repositories.DesignerProfileSnapshotRepository
}

func NewEventHandlers(repository repositories.DesignerProfileSnapshotRepository) EventHandlers {
	return EventHandlers{
		designerProfileSnapshotRepository: repository,
	}
}

func (e *EventHandlers) HandleKafkaEvent(message []byte) {
	var kafkaEvent KafkaEvent
	if err := json.Unmarshal(message, &kafkaEvent); err != nil {
		fmt.Printf("не удалось распрасить сообщение из kafka: %v", err)
	}

	switch kafkaEvent.eventType {
	case designerProfileCreated:
		profile, err := e.DecodeDomain(kafkaEvent.payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		e.designerProfileSnapshotRepository.CreateDesignerProfileSnapshot(profile)
	case designerProfileUpdated:
		profile, err := e.DecodeDomain(kafkaEvent.payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		e.designerProfileSnapshotRepository.UpdateDesignerProfileSnapshot(profile)
	case designerProfileDeleted:
		profile, err := e.DecodeDomain(kafkaEvent.payload)
		if err != nil {
			fmt.Printf(err.Error())
		}
		e.designerProfileSnapshotRepository.DeleteDesignerProfileSnapshot(profile.ID)
	}
}

func (e *EventHandlers) DecodeDomain(message json.RawMessage) (dto.DesignerProfileSnapshot, error) {
	var designerProfileSnapshot dto.DesignerProfileSnapshot
	if err := json.Unmarshal(message, &designerProfileSnapshot); err != nil {
		return dto.DesignerProfileSnapshot{}, fmt.Errorf("не удалось распрасить сущность внутри сообщения из kafka: %w", err)
	}

	return designerProfileSnapshot, nil
}
