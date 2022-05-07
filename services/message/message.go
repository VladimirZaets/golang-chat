package message

import (
	"time"
)

type Message struct {
	Name string    `json:"name"`
	Time time.Time `json:"created_at"`
	Data string    `json:"data"`
}

func (m *Message) Message() string {
	return m.Data
}
