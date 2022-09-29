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

func chatSocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("socket request")
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

// chatSocket struct
type chatSocket struct {
	con  *websocket.Conn
	mode int
	name string
}

func (i *chatSocket) broadcast(str string) {
	for _, g := range savedChatSockets {

		if g == i {
			// no send message to himself
			continue
		}

		if g.mode == 1 {
			// no send message to connected user before user write his name
			continue
		}
		g.writeMsg(i.name, str)
	}
}

func (i *chatSocket) read() {
	_, b, er := i.con.ReadMessage()
	if er != nil {
		panic(er)
	}
	fmt.Println(i.name + " " + string(b))
	fmt.Println(i.mode)

	if i.mode == 1 {
		i.name = string(b)
		i.writeMsg("Admin", "Welcome "+i.name+"!")
		i.mode = 2 // real msg mode
		return
	}
	i.broadcast(string(b))
	fmt.Println(i.name + " " + string(b))
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
			fmt.Println("thread socketreader finish")
		}()

		for {
			i.read()
		}
	}()
}
