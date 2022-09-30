package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

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

func chatSocketCreate(w http.ResponseWriter, r *http.Request) {
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

func (i *chatSocket) broadcast(str string) {
	for _, currentChatSocket := range savedChatSockets {
		if currentChatSocket == i {
			// users cannot send messages to themselves
			continue
		}
		if currentChatSocket.mode == 1 {
			// message cannot be sent until username is given
			continue
		}
		currentChatSocket.writeMsg(i.username, str)
	}
}

func (i *chatSocket) read() {
	_, b, er := i.con.ReadMessage()
	if er != nil {
		panic(er)
	}
	fmt.Println(i.username + " " + string(b))
	fmt.Println(i.mode)

	if i.mode == 1 {
		i.username = string(b)
		i.writeMsg("Admin", "Welcome "+i.username+"!")
		i.mode = 2 // real msg mode
		return
	}
	i.broadcast(string(b))
	fmt.Println(i.username + " " + string(b))
}

func (i *chatSocket) writeMsg(name string, str string) {
	i.con.WriteMessage(websocket.TextMessage, []byte("<b>"+name+": </b>"+str))
}

func (i *chatSocket) startThread() {
	i.writeMsg("Admin", "Please enter your username.")
	i.mode = 1 // mode 1 get user name

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
