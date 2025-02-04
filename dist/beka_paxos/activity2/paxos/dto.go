package paxos

import "encoding/json"

type QueueMessage struct {
	Type string          `json:"Type"`
	Body json.RawMessage `json:"Body"`
}
