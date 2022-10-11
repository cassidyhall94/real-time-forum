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
				fmt.Println("recovered panic:", err)
			}
			fmt.Println("pollContentWS finished")
		}()

		for {
			_, b, err := i.con.ReadMessage()
			if err != nil {
				panic(err)
			}
			tpl, err := template.ParseGlob("templates/*")
			if err != nil {
				panic(err)
			}
			w, err := i.con.NextWriter(websocket.TextMessage)
			if err != nil {
				panic(err)
			}
			switch string(b) {
			case "home":
				if err := tpl.ExecuteTemplate(w, "home.template", nil); err != nil {
					panic(fmt.Errorf("Home ExecuteTemplate error: %w", err))
				}
			case "posts":
				if err := tpl.ExecuteTemplate(w, "posts.template", nil); err != nil {
					panic(fmt.Errorf("Posts ExecuteTemplate error: %w", err))
				}
			case "login":
				if err := tpl.ExecuteTemplate(w, "login.template", nil); err != nil {
					panic(fmt.Errorf("Login ExecuteTemplate error: %w", err))
				}
			default:
				panic(fmt.Errorf("template %s not found", string(b)))
			}
			if err := w.Close(); err != nil {
				panic(err)
			}
		}
	}()
}
