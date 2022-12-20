package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/authentication"
	"real-time-forum/pkg/database"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
)

type SocketMessage struct {
	Type messageType `json:"type,omitempty"`
}

type socket struct {
	con     *websocket.Conn
	user    *database.User
	t       messageType
	uuid    uuid.UUID
	created time.Time
}

type socketTimeoutError struct {
	Message string
}

func (s *socketTimeoutError) Error() string {
	if s.Message == "" {
		return "socket timed out mate"
	}
	return s.Message
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	savedSockets    = make([]*socket, 0)
	timeoutDuration = 1 * time.Hour
)

func SocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Socket Request on " + r.RequestURI)
	if !authentication.RequestIsLoggedIn(r) {
		fmt.Println("unauthorised socket creation request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u, err := authentication.GetUserFromRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ptrSocket := &socket{
		con:     con,
		user:    u,
		uuid:    uuid.NewV4(),
		created: time.Now(),
	}
	defer func() {
		for i, so := range savedSockets {
			if so.uuid == ptrSocket.uuid {
				savedSockets = removeFromSlice(savedSockets, i)
				fmt.Printf("removed '%s' socket (ID: '%s') belonging to '%s (ID: %s)'\n", ptrSocket.t.String(), ptrSocket.uuid, ptrSocket.user.Nickname, ptrSocket.user.ID)
			}
		}
	}()

	// add new case here when added to main.go for handlers
	switch r.RequestURI {
	case "/content":
		ptrSocket.t = content
		// loads the home page (which contains the posts form)

		if err := OnContentConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	case "/post":
		ptrSocket.t = post
		// loads the saved posts on window load
		if err := OnPostsConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	case "/chat":
		ptrSocket.t = chat
	case "/presence":
		ptrSocket.t = presence
		// loads the presence list on window load
		if err := OnPresenceConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	default:
		ptrSocket.t = unknown
	}

	savedSockets = append(savedSockets, ptrSocket)
	w.Write([]byte{})
	if err := ptrSocket.pollSocket(); err != nil {
		fmt.Println(err)
		return
	}
}
func (s *socket) pollSocket() error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		for {
			if s.IsTimedOut() {
				return &socketTimeoutError{}
			}
			b, err := s.read()
			if err != nil {
				return fmt.Errorf("error reading from socket '%s': %w", s.uuid, err)
			} else if b == nil {
				return fmt.Errorf("socket '%s' closed", s.uuid)
			}
			sm := &SocketMessage{}
			if err := json.Unmarshal(b, sm); err != nil {
				return fmt.Errorf("error unmarshalling message in socket '%s': %w", s.uuid, err)
			}
			switch sm.Type {
			case chat:
				m := &ChatMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					return fmt.Errorf("error unmarshalling chatMessage in socket '%s': %w", s.uuid, err)
				}
				if err := m.Handle(s); err != nil {
					return fmt.Errorf("error handling chatMessage in socket '%s': %w", s.uuid, err)
				}
			case post:
				m := &PostMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					return fmt.Errorf("error unmarshalling postMessage in socket '%s': %w", s.uuid, err)
				}
				if err := m.Handle(s); err != nil {
					return fmt.Errorf("error handling postMessage in socket '%s': %w", s.uuid, err)
				}
			case content:
				m := &ContentMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					return fmt.Errorf("error unmarshalling contentMessage in socket '%s': %w", s.uuid, err)
				}
				if err := m.Handle(s); err != nil {
					return fmt.Errorf("error handling contentMessage in socket '%s': %w", s.uuid, err)
				}
			case presence:
				m := &PresenceMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					return fmt.Errorf("error unmarshalling presenceMessage in socket '%s': %w", s.uuid, err)
				}
				if err := m.Handle(s); err != nil {
					return fmt.Errorf("error handling presenceMessage in socket '%s': %w", s.uuid, err)
				}
			default:
				panic(fmt.Errorf("unable to determine message type for '%s'", string(b)))
			}
		}
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("encountered error when polling socket %+v: %w", *s, err)
	}

	return nil
}

func (s *socket) read() ([]byte, error) {
	_, b, err := s.con.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return nil, nil
		}
		log.Print(b)
		return nil, fmt.Errorf("unable to read message from socket, got: '%s', %w", string(b), err)
	}
	return b, nil
}

func removeFromSlice[T any](slice []T, index int) []T {
	// Check that the index is within the bounds of the slice
	if index < 0 || index >= len(slice) {
		return slice
	}

	// Remove the value at the specified index by replacing it with the
	// last value in the slice and then slicing off the last element
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func (s *socket) IsTimedOut() bool {
	if s.created.Add(timeoutDuration).After(time.Now()) {
		return false
	}
	return true
}
