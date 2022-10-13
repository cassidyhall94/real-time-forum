package websockets

import (
	"fmt"
)

type PostMessage struct {
	Type       messageType `json:"type,omitempty"`
	Header     string      `json:"text,omitempty"`
	Body       string      `json:"body,omitempty"`
	Categories []string    `json:"categories,omitempty"`
	Timestamp  string      `json:"timestamp,omitempty"`
	Username   string      `json:"username,omitempty"`
}

func (m PostMessage) Handle(s *socket) error {
	return nil
}

func (m *PostMessage) Broadcast() error {
	for _, s := range savedSockets {
		if s.t == m.Type {
			if err := s.con.WriteJSON(m); err != nil {
				return fmt.Errorf("unable to send message: %w", err)
			}
		}
	}
	return nil
}
