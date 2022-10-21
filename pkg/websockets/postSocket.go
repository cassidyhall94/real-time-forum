package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
)

type PostMessage struct {
	Type      messageType     `json:"type,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	Posts     []database.Post `json:"posts,omitempty"`
}

func (m PostMessage) Handle(s *socket) error {
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

func GetPostsOnConnect() ([]database.Post, error) {
	posts, err := database.GetPosts()
	if err != nil {
		return posts, fmt.Errorf("GetPostsOnConnect (GetPosts) error: %+v\n", err)
	}

	for _, post := range posts {
		posts = append(posts, database.Post{
			PostID:     post.PostID,
			Username:   post.Username,
			Title:      post.Title,
			Categories: post.Categories,
			Body:       post.Body,
			// Comments: ,
		})
	}
	return posts, nil
}

func OnPostsConnect(s *socket) error {
	posts, err := GetPostsOnConnect()
	if err != nil {
		return fmt.Errorf("OnPostsConnect (GetPostsOnConnect) error: %+v\n", err)
	}

	for _, post := range posts {
		posts = append(posts, database.Post{
			PostID:     post.PostID,
			Username:   post.Username,
			Title:      post.Title,
			Categories: post.Categories,
			Body:       post.Body,
			// Comments: ,
		})
	}
	c := &PostMessage{
		Type:      post,
		Timestamp: "",
		Posts:     posts,
	}
	return c.Broadcast(s)
}
