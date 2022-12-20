package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"sort"
	"time"
)

type PresenceMessage struct {
	Type      messageType          `json:"type"`
	Timestamp string               `json:"timestamp"`
	Presences []*database.Presence `json:"presences"`
}

func (m *PresenceMessage) Broadcast(s *socket) error {
	if s.IsTimedOut() {
		return &socketTimeoutError{}
	}

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

func OnPresenceConnect(s *socket) error {
	for {
		presences, err := database.GetPresencesForUser(s.user)
		if err != nil {
			return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v\n", err)
		}
		sort.SliceStable(presences, func(p, q int) bool {
			if presences[p].Online {
				return true
			} else if presences[p].LastContactedTime.Before(presences[q].LastContactedTime) {
				return true
			} else if presences[p].User.Nickname < presences[q].User.Nickname {
				return true
			} else {
				return false
			}
		})
		c := &PresenceMessage{
			Type:      presence,
			Timestamp: time.Now().String(),
			Presences: presences,
		}
		if err := c.Broadcast(s); err != nil {
			return fmt.Errorf("OnPresenceConnect (c.Broadcast) error: %w", err)
		}
		time.Sleep(1 * time.Second)
	}
}