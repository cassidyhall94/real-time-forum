package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"time"
)

type PostMessage struct {
	Type      messageType     `json:"type,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	Posts     []database.Post `json:"posts,omitempty"`
}

func (m PostMessage) Handle(s *socket) error {
	for _, post := range m.Posts {
		CreatePost(post)
	}
	return m.Broadcast(s)
}

func (m *PostMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (post) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send post message down ws of type %s", s.t.String())
	}
	return nil
}

func OnPostsConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	posts, err := database.GetPosts()
	if err != nil {
		return fmt.Errorf("OnPostsConnect (GetPosts) error: %+v\n", err)
	}
	c := &PostMessage{
		Type:      post,
		Timestamp: "",
		Posts:     posts,
	}
	return c.Broadcast(s)
}

func CreatePost(post database.Post) error {
	fmt.Println(post)
	stmt, err := database.DB.Prepare("INSERT INTO posts (postID, username, title, categories, body) VALUES (uuid.NewV4(), ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(post.PostID, post.Username, post.Title, post.Body, post.Categories)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return nil
}
