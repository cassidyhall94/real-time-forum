package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
)

// PresenceMessage contains some meta data about PresenceMessages and contains a []Presence (ID is contained in the database users table)
type PresenceMessage struct {
	Type      messageType `json:"type"`
	Timestamp string      `json:"timestamp"`
	Presences []Presence  `json:"presences"`
}

type Presence struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Online            bool   `json:"online"`
	LastContactedTime string `json:"last_contacted_time"`
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

func GetPresences() ([]Presence, error) {
	presences := []Presence{}

	users, err := database.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("GetUsers (getPresences) error: %+v\n", err)
	}
	for _, user := range users {
		presences = append(presences, Presence{
			ID:       user.ID,
			Username: user.Username,
			// Online:            bool,
			// LastContactedTime: created,
		})
	}
	return presences, nil
}

func OnPresenceConnect(s *socket) error {
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v\n", err)
	}
	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: presences,
		// Presences: []Presence{
		// {
		// 	ID:                "id1",
		// 	Username:          "user1",
		// 	Online:            true,
		// 	LastContactedTime: "213243532",
		// },
		// {
		// 	ID:                "id2",
		// 	Username:          "user2",
		// 	Online:            true,
		// 	LastContactedTime: "2132432",
		// },
		// {
		// 	ID:                "id3",
		// 	Username:          "user3",
		// 	Online:            false,
		// 	LastContactedTime: "2133532",
		// },
		// },
	}
	return c.Broadcast(s)
}

// func (data *Forum) GetSessions() ([]Session, error) {
// 	session := []Session{}
// 	rows, err := data.DB.Query(`SELECT * FROM session`)
// 	if err != nil {
// 		return session, fmt.Errorf("GetSession DB Query error: %+v\n", err)
// 	}
// 	var session_token string
// 	var uName string
// 	var exTime string
// 	for rows.Next() {
// 		err := rows.Scan(&session_token, &uName, &exTime)
// 		if err != nil {
// 			return nil, fmt.Errorf("GetSession rows.Scan error: %+v\n", err)
// 		}
// 		// time.Format("01-02-2006 15:04")
// 		convTime, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", exTime)
// 		if err != nil {
// 			return nil, fmt.Errorf("GetSession time.Parse(exTime) error: %+v\n", err)
// 		}
// 		session = append(session, Session{
// 			SessionID: session_token,
// 			Username:  uName,
// 			Expiry:    convTime,
// 		})
// 	}
// 	return session, nil
// }
