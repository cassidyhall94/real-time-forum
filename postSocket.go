package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var savedPostSockets []*postSocket

// chatSocket struct
type postSocket struct {
	con      *websocket.Conn
	mode     int
	username string
}

func postSocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Socket Request")
	if savedPostSockets == nil {
		savedPostSockets = make([]*postSocket, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		r.Body.Close()
	}()
	con, _ := upgrader.Upgrade(w, r, nil)
	ptrPostSocket := &postSocket{
		con: con,
	}

	savedPostSockets = append(savedPostSockets, ptrPostSocket)
	ptrPostSocket.startThread()
}

func (i *postSocket) broadcast(str string) {
	for _, currentPostSocket := range savedPostSockets {
		if currentPostSocket == i {
			// users cannot send messages to themselves
			continue
		}
		if currentPostSocket.mode == 1 {
			// message cannot be sent until username is given
			continue
		}
		currentPostSocket.writeMsg(i.username, str)
	}
}

func (i *postSocket) read() {
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
	// fmt.Println(i.username + " " + string(b))
}

func (i *postSocket) writeMsg(name string, str string) {
	i.con.WriteMessage(websocket.TextMessage, []byte("<b>"+dateTime+" </b>"+"<br>"+"<b>"+name+": </b>"+str))
}

func (i *postSocket) startThread() {
	i.writeMsg("Admin", "Please enter your username.")
	i.mode = 1 // mode 1 get user name

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Post thread finished")
		}()

		for {
			i.read()
		}
	}()
}
