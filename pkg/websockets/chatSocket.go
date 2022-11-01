package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
)

type ChatMessage struct {
	Type          messageType              `json:"type,omitempty"`
	Timestamp     string                   `json:"timestamp,omitempty"`
	Conversations []*database.Conversation `json:"conversations"`
}

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
		for _, convo := range conversations {
			for _, chat := range convo.Chats {
				fmt.Println("handle: ", chat)
			}
		}
		return c.Broadcast(s)
	}
	for _, convo := range m.Conversations {
		if convo.ConvoID == "" {
			newConvoID, err := CreateConversation(convo)
			convo.ConvoID = newConvoID
			if err != nil {
				return fmt.Errorf("ChatSocket Handle (CreateConversation) error: %w", err)
			}

		}
		for _, chat := range convo.Chats {
			if chat.ChatID == "" {
				newChatID, err := CreateChat(chat)
				chat.ChatID = newChatID
				if err != nil {
					return fmt.Errorf("ChatSocket Handle (CreateChat) error: %w", err)
				}
			}
		}
		// the id for the chat/convo now exists in the DB as a result of Create*() but the ID is not written back into m
		// this means that the message that's send in m.Broadcast is missing these newly created ID's
		return m.Broadcast(s)
	}
	return nil
}

func CreateChat(chat database.Chat) (string, error) {
	stmt, err := database.DB.Prepare("INSERT INTO chats (convoID, chatID, sender, date, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("CreateChat DB Prepare error: %+v\n", err)
	}
	if chat.ChatID == "" {
		chat.ChatID = uuid.NewV4().String()
	}

	// TODO: remove placeholder nickname once login/sessions are working
	if chat.Sender.ID == "" {
		//this is foo's userID in the database
		chat.Sender.ID = "6d01e668-2642-4e55-af73-46f057b731f9"
	}

	_, err = stmt.Exec(chat.ConvoID, chat.ChatID, chat.Sender.ID, chat.Date, chat.Body)
	if err != nil {
		return "", fmt.Errorf("CreateChat Exec error: %+v\n", err)
	}
	return chat.ChatID, err
}

func CreateConversation(conversations *database.Conversation) (string, error) {
	stmt, err := database.DB.Prepare("INSERT INTO conversations (convoID, participants) VALUES (?, ?);")
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("CreateConversations DB Prepare error: %+v\n", err)
	}
	if conversations.ConvoID == "" {
		conversations.ConvoID = uuid.NewV4().String()
	}
	for _, participant := range conversations.Participants {
		_, err = stmt.Exec(conversations.ConvoID, participant.ID)
		if err != nil {
			return "", fmt.Errorf("CreateConversations Exec error: %+v\n", err)
		}
	}
	return conversations.ConvoID, err
}
