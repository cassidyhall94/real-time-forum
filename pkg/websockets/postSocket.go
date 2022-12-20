package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"time"
)

type PostMessage struct {
	Type      messageType      `json:"type,omitempty"`
	Timestamp string           `json:"timestamp,omitempty"`
	Posts     []*database.Post `json:"posts"`
}

func (m PostMessage) Handle(s *socket) error {
	if len(m.Posts) == 0 {
		p, err := database.GetPopulatedPosts()
		if err != nil {
			return err
		}
		c := &PostMessage{
			Type:  post,
			Posts: p,
		}
		return c.Broadcast(s)
	}
	for _, pt := range m.Posts {
		if pt.PostID == "" {
			newPost, err := database.CreatePost(pt)
			if err != nil {
				return fmt.Errorf("PostSocket Handle (CreatePost) error: %w", err)
			}
			pt = newPost
		}
		for _, comment := range pt.Comments {
			if comment.CommentID == "" {
				newComment, err := database.CreateComment(comment)
				if err != nil {
					return fmt.Errorf("PostSocket Handle (CreateComment) error: %w", err)
				}
				comment = newComment
			}
		}
		return m.Broadcast(s)
	}
	return nil
}

func (m *PostMessage) Broadcast(s *socket) error {
	if s.IsTimedOut() {
		return &socketTimeoutError{}
	}

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
	p, err := database.GetPopulatedPosts()
	if err != nil {
		return err
	}
	c := &PostMessage{
		Type:      post,
		Timestamp: time.Now().String(),
		Posts:     p,
	}
	return c.Broadcast(s)
}
