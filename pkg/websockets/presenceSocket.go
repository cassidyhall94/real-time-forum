package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"sort"
	"time"
)

type PresenceMessage struct {
	Type      messageType         `json:"type"`
	Timestamp string              `json:"timestamp"`
	Presences []database.Presence `json:"presences"`
}

func (m *PresenceMessage) Broadcast(s *socket) error {
	// fmt.Println("hi broadcast")
	if s.t == m.Type {
		m.Presences, _ = GetPresences()
		// fmt.Println("GetPresences()", m.Presences)
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (presence) message: %w", err)
		}
		// OnPresenceConnect(s)
	} else {
		return fmt.Errorf("cannot send presence message down ws of type %s", s.t.String())
	}
	return nil
}

func (m *PresenceMessage) Handle(s *socket) error {
	return m.Broadcast(s)
}

func GetPresences() ([]database.Presence, error) {
	// fmt.Println("hi get presences")
	presences := []database.Presence{}
	users, err := database.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("GetUsers (getPresences) error: %+v\n", err)
	}
	sort.SliceStable(users[:], func(i, j int) bool {
		return users[i].Nickname < users[j].Nickname
	})
	for _, user := range users {
		if user.LoggedIn == "true" {
			// fmt.Println("testing", user.Nickname)
			presences = append(presences, database.Presence{
				ID:       user.ID,
				Nickname: user.Nickname,
				Online:   user.LoggedIn,
				// LastContactedTime: created,
				//TODO CHAT: create a getLastContactedTime() func for getting the timestamp of the last message sent to clickedParticipantID by currentUserID to organise presence list and chat messages
			})

		}
	}
	// fmt.Println("get presences", presences)
	return presences, nil
}

func OnPresenceConnect(s *socket) error {
	// fmt.Println("hguefgsuyuufygsyegfgefuag8efga8gef87aef87aef78ef")

	time.Sleep(1 * time.Second)
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v\n", err)
	}
	c := &PresenceMessage{
		Type:      presence,
		Timestamp: time.Now().String(),
		Presences: presences,
	}
	return c.Broadcast(s)
}
