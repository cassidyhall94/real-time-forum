package websockets

import "fmt"

type ChatMessage struct {
	Type      messageType `json:"type,omitempty"`
	Text      string      `json:"text,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Username  string      `json:"username,omitempty"`
}

func (m *ChatMessage) Broadcast() error {
	for _, s := range savedSockets {
		if s.t == m.Type {
			if err := s.con.WriteJSON(m); err != nil {
				return fmt.Errorf("unable to send message: %w", err)
			}
		}
	}
	return nil
}

func (m *ChatMessage) Handle(s *socket) error {
	return nil
}
