package floodcontrol

import "time"

type Payload struct {
	LastCallAt time.Time `json:"last_call_at"`
	CallCount  int       `json:"call_count"`
}
