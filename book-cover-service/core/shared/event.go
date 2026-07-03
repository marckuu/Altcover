package shared

import "encoding/json"

const (
	designerProfileCreated = iota
	designerProfileUpdated
	designerProfileDeleted
)

type KafkaEvent struct {
	EventType int
	Payload   json.RawMessage
}
