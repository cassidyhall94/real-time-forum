package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
)

type ChatMessage struct {
	Type          messageType              `json:"type,omitempty"`
	Timestamp     string                   `json:"timestamp,omitempty"`
	Conversations []*database.Conversation `json:"chats"`
}

// type Conversation struct {
// 	ChatID       string   `json:"chat_id,omitempty"`
// 	Participants []string `json:"participants"`
// 	Chats        []Chat   `json:"chats,omitempty"`
// }

// type Chat struct {
// 	ChatID   string `json:"chat_id"`
// 	Sender   string `json:"sender"`
// 	Receiver string `json:"receiver"`
// 	Date     string `json:"date,omitempty"`
// 	Body     string `json:"body,omitempty"`
// }

func (m *ChatMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (chat) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send chat message down ws of type %s", s.t.String())
	}
	return nil
}

func (m *ChatMessage) Handle(s *socket) error {
	if len(m.Conversations) == 0 {
		conversations, err := database.GetPopulatedConversations()
		if err != nil {
			return err
		}

		c := &ChatMessage{
			Type:          chat,
			Conversations: conversations,
		}

		return c.Broadcast(s)
	}
	for _, chats := range m.Conversations {
		if chats.ChatID == "" {
			if err := CreateConversations(chats); err != nil {
				return fmt.Errorf("ChatSocket Handle (CreateConversations) error: %w", err)
			}
		}
		for _, chat := range chats.Chats {
			if chat.ChatID == "" {
				if err := CreateChat(chat); err != nil {
					return fmt.Errorf("ChatSocket Handle (CreateChat) error: %w", err)
				}
			}
		}
	}
	return nil
}

func CreateChat(chat database.Chat) error {
	stmt, err := database.DB.Prepare("INSERT INTO chat (chatID, sender, receiver, date, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateChat DB Prepare error: %+v\n", err)
	}
	if chat.ChatID == "" {
		chat.ChatID = uuid.NewV4().String()
	}

	// TODO: remove placeholder nickname once login/sessions are working
	if chat.Sender == "" {
		chat.Sender = "Cassidy"
	}

	if chat.Receiver == "" {
		chat.Sender = "Jeff"
	}

	_, err = stmt.Exec(chat.ChatID, chat.Sender, chat.Receiver, chat.Date, chat.Body)
	if err != nil {
		return fmt.Errorf("CreateChat Exec error: %+v\n", err)
	}
	return nil
}

func CreateConversations(conversations *database.Conversation) error {
	stmt, err := database.DB.Prepare("INSERT INTO conversations (chatID, participants) VALUES (?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateConversations DB Prepare error: %+v\n", err)
	}
	if conversations.ChatID == "" {
		conversations.ChatID = uuid.NewV4().String()
	}

	_, err = stmt.Exec(conversations.ChatID, conversations.Participants)
	if err != nil {
		return fmt.Errorf("CreateConversations Exec error: %+v\n", err)
	}
	return nil
}
