package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	auth "real-time-forum/pkg/authentication"
	"time"

	"github.com/gorilla/websocket"
)

type SocketMessage struct {
	Type    messageType `json:"type,omitempty"`
}

type socket struct {
	con      *websocket.Conn
	username string
	t        messageType
}

var savedSockets []*socket

var (
	t        = time.Now()
	dateTime = t.Format("1/2/2006, 3:04:05 PM")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func SocketCreate(w http.ResponseWriter, r *http.Request) {
	c1, err1 := r.Cookie("1st-cookie")
	if err1 == nil && !auth.Person.Accesslevel {
		// first home page access
		c1.MaxAge = -1
		http.SetCookie(w, c1)
	}
	_, err := r.Cookie("1st-cookie")
	if err != nil && auth.Person.Accesslevel {
		// logged in and on 2nd browser
		auth.Person.CookieChecker = false
	} else if err == nil && auth.Person.Accesslevel {
		// Original browser and logged in
		auth.Person.CookieChecker = true
	} else {
		// not logged in yet
		auth.Person.CookieChecker = false
	}

	fmt.Println("Socket Request on " + r.RequestURI)
	if savedSockets == nil {
		savedSockets = make([]*socket, 0)
	}

	con, _ := upgrader.Upgrade(w, r, nil)
	ptrSocket := &socket{
		con: con,
	}
	// add new case here when added to main.go for handlers
	switch r.RequestURI {
	case "/content":
		ptrSocket.t = content

		// before doing anything else, send the splash(home) page
		if err := OnContentConnect(ptrSocket); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/post":
		ptrSocket.t = post
	case "/chat":
		ptrSocket.t = chat
	default:
		ptrSocket.t = unknown
	}

	savedSockets = append(savedSockets, ptrSocket)
	ptrSocket.pollSocket()
}

func (s *socket) pollSocket() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("recovered panic in %s socket: %+v\n", s.t.String(), err)
			}
			fmt.Println(s.t.String() + " socket closed")
		}()

		for {
			b, err := s.read()
			if err != nil {
				panic(err)
			}
			sm := &SocketMessage{}
			if err := json.Unmarshal(b, sm); err != nil {
				panic(err)
			}
			switch sm.Type {
			case chat:
				m := &ChatMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			case post:
				m := &PostMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			case content:
				m := &ContentMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			default:
				panic(fmt.Errorf("unable to determine message type for '%s'", string(b)))
			}
		}
	}()
}

func (s *socket) read() ([]byte, error) {
	_, b, err := s.con.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("unable to read message from socket, got: '%s', %w", string(b), err)
	}
	return b, nil
}