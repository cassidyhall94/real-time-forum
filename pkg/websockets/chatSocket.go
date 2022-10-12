package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	auth "real-time-forum/pkg/authentication"
	"time"

	"github.com/gorilla/websocket"
)

var t = time.Now()
var dateTime = t.Format("1/2/2006, 3:04:05 PM")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// eventually === map[string (username?)]*chatSocket
var savedChatSockets []*chatSocket

// chatSocket struct
type chatSocket struct {
	con      *websocket.Conn
	username string
}

type ChatMessage struct {
	PresenceList []string
	Text         string
	Timestamp    string
	Username     string
}

func ChatSocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Chat Socket Request")
	if savedChatSockets == nil {
		savedChatSockets = make([]*chatSocket, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		r.Body.Close()
	}()
	con, _ := upgrader.Upgrade(w, r, nil)
	ptrChatSocket := &chatSocket{
		con: con,
	}

	savedChatSockets = append(savedChatSockets, ptrChatSocket)
	ptrChatSocket.startThread()
}

func (i *chatSocket) broadcast(c *ChatMessage) error {
	for range savedChatSockets {
		// when savedChatSockets is turned into a map[username]socket, this will need to be changed to filter for specific usernames (DMing)
		if err := i.con.WriteJSON(c); err != nil {
			return fmt.Errorf("unable to send chat message: %w", err)
		}
	}
	return nil
}

func (i *chatSocket) read() {
	_, b, er := i.con.ReadMessage()
	if er != nil {
		panic(er)
	}
	fmt.Println(string(b))
	c := &ChatMessage{}
	if err := json.Unmarshal(b, c); err != nil {
		panic(err)
	}
	c.Username = auth.Person.Username
	i.broadcast(c)
}

func (i *chatSocket) startThread() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Chat thread finished")
		}()

		for {
			i.read()
		}
	}()
}
