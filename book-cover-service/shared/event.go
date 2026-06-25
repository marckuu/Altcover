package shared

import "encoding/json"

const (
	designerProfileCreated = iota
	designerProfileUpdated
	designerProfileDeleted
)

type KafkaEvent struct {
	eventType int
	payload   json.RawMessage
}
