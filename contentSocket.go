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

var contentUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var savedContentSockets []*contentSocket



func contentSocketCreate(w http.ResponseWriter, r *http.Request) {

		c1, err1 := r.Cookie("1st-cookie")
	if err1 == nil && !Person.Accesslevel {
		// first home page access 
		c1.MaxAge = -1
		http.SetCookie(w, c1)
	}
	_, err := r.Cookie("1st-cookie")
	if err != nil && Person.Accesslevel {
		// logged in and on 2nd browser
		Person.CookieChecker = false
	} else if err == nil && Person.Accesslevel {
		// Original browser and logged in
		Person.CookieChecker = true
	} else {
		// not logged in yet
	Person.CookieChecker = false
	}

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
					fmt.Printf("Home ExecuteTemplate error: %+v\n", err)
					return
				}
			case "reg-log":
				if err := tpl.ExecuteTemplate(w, "reg-log.template", nil); err != nil {
					fmt.Printf("Reg-Log ExecuteTemplate error: %+v\n", err)
					return
				}
			}
		}
	}()
}
