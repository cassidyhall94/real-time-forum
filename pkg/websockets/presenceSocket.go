package websockets

import (
	"fmt"
)

// PresenceMessage contains some meta data about PresenceMessages and contains a map[ID]Presence (ID is contained in the database users table)
type PresenceMessage struct {
	Type      messageType         `json:"type"`
	Timestamp string              `json:"timestamp"`
	Presences map[string]Presence `json:"presences"`
}

type Presence struct {
	Username          string `json:"username"`
	Online            bool   `json:"online"`
	LastContactedTime string `json:"last_contacted_time"`
}

func OnPresenceConnect(s *socket) error {

	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: map[string]Presence{
			"id1": {
				Username: "user1",
				Online: true,
				LastContactedTime: "213243532",
			},
			"id2": {
				Username: "user2",
				Online: true,
				LastContactedTime: "2132432",
			},
			"id3": {
				Username: "user3",
				Online: false,
				LastContactedTime: "2133532",
			},
		},
	}

	return c.Broadcast(s)
}

func (m *PresenceMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (presence) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send presence message down ws of type %s", s.t.String())
	}
	return nil
}

func (m *PresenceMessage) Handle(s *socket) error {
	return m.Broadcast(s)
}
