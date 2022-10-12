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

var savedChatSockets []*chatSocket

// chatSocket struct
type chatSocket struct {
	con      *websocket.Conn
	mode     int
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
	j, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("unable to marshal chat message: %w", err)
	}
	for _, currentChatSocket := range savedChatSockets {
		if currentChatSocket == i {
			// users cannot send messages to themselves
			continue
		}
		if currentChatSocket.mode == 1 {
			// message cannot be sent until username is given
			continue
		}
		if err := i.con.WriteMessage(websocket.TextMessage, j); err != nil {
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
	if i.mode == 1 {
		c.Username = auth.Person.Username
		i.mode = 2 // real msg mode
		return
	}
	i.broadcast(c)
}

func (i *chatSocket) writeMsg(name string, str string) {

	i.con.WriteMessage(websocket.TextMessage, []byte("<b>"+dateTime+" </b>"+"<br>"+"<b>"+name+": </b>"+str))
}

func (i *chatSocket) startThread() {
	i.mode = 1
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
