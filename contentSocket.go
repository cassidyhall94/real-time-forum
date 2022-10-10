package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

var savedContentSocket *contentSocket

// contentSocket struct
type contentSocket struct {
	con *websocket.Conn
	// mode int
	template string
}

func contentSocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Content Socket Request")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
		r.Body.Close()
	}()
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	ptrContentSocket := &contentSocket{
		con: con,
	}

	savedContentSocket = ptrContentSocket
	ptrContentSocket.pollContentWS()
}

// pollContentWS starts listening on a websocket for messages
func (i *contentSocket) pollContentWS() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("pollContentWS finished")
		}()

		for {
			_, b, err := i.con.ReadMessage()
			if err != nil {
				fmt.Println(err)
			}
			tpl, err := template.ParseGlob("templates/*")
			if err != nil {
				fmt.Println(err)
			}
			w, err := i.con.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
			}
			switch string(b) {
			case "home":
				if err := tpl.ExecuteTemplate(w, "home.template", nil); err != nil {
					fmt.Printf("Home ExecuteTemplate error: %+v\n", err)
					return
				}
			case "posts":
				if err := tpl.ExecuteTemplate(w, "posts.template", nil); err != nil {
					fmt.Printf("Posts ExecuteTemplate error: %+v\n", err)
					return
				}
			case "login":
				if err := tpl.ExecuteTemplate(w, "login.template", nil); err != nil {
					fmt.Printf("Login ExecuteTemplate error: %+v\n", err)
					return
				}
			}
		}
	}()
}
