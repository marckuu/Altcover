package shared

import "encoding/json"

const (
	DesignerProfileCreated = iota
	DesignerProfileUpdated
	DesignerProfileDeleted
)

type KafkaEvent struct {
	EventType int
	Payload   json.RawMessage
}

func NewKafkaEvent(eventType int, payload json.RawMessage) KafkaEvent {
	return KafkaEvent{
		EventType: eventType,
		Payload:   payload,
	}
}
