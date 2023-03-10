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
	if s.IsTimedOut() {
		return &socketTimeoutError{}
	}

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
	// fmt.Printf("\nm: %+v\n", *m.Conversations[0])
	if len(m.Conversations) == 0 {
		conversations, err := database.GetPopulatedConversations(nil)
		if err != nil {
			return err
		}
		c := &ChatMessage{
			Type:          chat,
			Conversations: conversations,
		}
		return c.Broadcast(s)
	}

	allConversations, err := database.GetConversations()
	if err != nil {
		return err
	}

	for i, convo := range m.Conversations {
		// creates a new conversation if the convoID is missing
		if len(convo.ConvoID) == 0 {
			id, err := database.GetConvoID(database.ParticipantsToIds(convo.Participants), allConversations)
			if err != nil {
				fmt.Printf("ChatSocket Handle (GetConvoID) error: %+v\n", err)
			}
			if err != nil || id == "" {
				newConvoID, err := CreateConversation(convo)
				if err != nil {
					return fmt.Errorf("ChatSocket Handle (CreateConversation) error: %w", err)
				}
				convo.ConvoID = newConvoID
			} else {
				convo.ConvoID = id
			}
		}
		for j, chat := range convo.Chats {
			// for new chats, the chat.ConvoID is given the conversation's convoID if it is missing
			if chat.ConvoID == "" {
				chat.ConvoID = convo.ConvoID
			}
			if chat.ChatID == "" {
				newChatID, err := CreateChat(chat)
				if err != nil {
					return fmt.Errorf("ChatSocket Handle (CreateChat) error: %w", err)
				}
				chat.ChatID = newChatID
			}
			convo.Chats[j] = chat
		}
		m.Conversations[i] = convo
	}
	c, err := database.GetPopulatedConversations(m.Conversations)
	if err != nil {
		return fmt.Errorf("ChatSocket Handle (GetPopulatedConversations) error: %w", err)
	}
	cm := &ChatMessage{
		Type:          chat,
		Conversations: c,
	}

	// b, err := json.Marshal(cm.Conversations)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("ChatSocket Handle (Convos): ", string(b))
	return cm.Broadcast(s)
}

func CreateChat(chat database.Chat) (string, error) {
	stmt, err := database.DB.Prepare("INSERT INTO chats (convoID, chatID, sender, date, body) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return "", fmt.Errorf("CreateChat DB Prepare error: %+v\n", err)
	}
	defer stmt.Close()
	if chat.ChatID == "" {
		chat.ChatID = uuid.NewV4().String()
	}
	fmt.Printf("sender userID: %+v\n", chat.Sender.ID)
	// TODO: remove placeholder nickname once login/sessions are working
	if chat.Sender.ID == "" {
		fmt.Printf("sender userID is blank, inserting foo's userID")
		//this is foo's userID in the database
		// chat.Sender.ID = "6d01e668-2642-4e55-af73-46f057b731f9"
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
