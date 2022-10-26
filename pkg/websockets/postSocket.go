package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PostMessage struct {
	Type      messageType      `json:"type,omitempty"`
	Timestamp string           `json:"timestamp,omitempty"`
	Posts     []*database.Post `json:"posts"`
}

// TODO: add code for handling comments and attaching to post
func (m PostMessage) Handle(s *socket) error {
	// create new post
	// if any posts in m.Posts is id == "" then make a new post
	// create new comment for post
	// if a post contains a comment with id == "" then create a comment for that post
	// if len(posts) == 0 then return all posts
	// else range over posts, get comments for post and override post.comments then return m
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

// TODO: add timestamp
func OnPostsConnect(s *socket) error {
	time.Sleep(1 * time.Second)

	p, err := database.GetPopulatedPosts()
	if err != nil {
		return err
	}

	c := &PostMessage{
		Type: post,
		Posts: p,
	}

	return c.Broadcast(s)
}

func CreatePost(post *database.Post) error {
	stmt, err := database.DB.Prepare("INSERT INTO posts (postID, username, title, categories, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	if post.PostID == "" {
		post.PostID = uuid.NewV4().String()
	}

	// TODO: remove placeholder username once login/sessions are working
	if post.Username == "" {
		post.Username = "Cassidy"
	}

	_, err = stmt.Exec(post.PostID, post.Username, post.Title, post.Categories, post.Body)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return nil
}

func CreateComment(comment database.Comment) error {
	stmt, err := database.DB.Prepare("INSERT INTO comments (commentID, postID, username, body) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	if comment.CommentID == "" {
		comment.CommentID = uuid.NewV4().String()
	}

	// TODO: remove placeholder username once login/sessions are working
	if comment.Username == "" {
		comment.Username = "Cassidy"
	}

	_, err = stmt.Exec(comment.CommentID, comment.PostID, comment.Username, comment.Body)
	if err != nil {
		return fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	return nil
}
