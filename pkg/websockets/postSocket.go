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
	// if len(posts) == 0 then return all posts
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
		// create new post
		// if any posts in m.Posts is id == "" then make a new post
		// else range over posts, get comments for post and override post.comments then return m
		for _, post := range m.Posts {
			if post.PostID == "" {
				if err := CreatePost(post); err != nil {
					return fmt.Errorf("PostSocket Handle (CreatePost) error: %w", err)
				}
			}
			// create new comment for post
			// if a post contains a comment with id == "" then create a comment for that post
			for _, comment := range post.Comments {
				if comment.CommentID == "" {
					if err := CreateComment(comment); err != nil {
						return fmt.Errorf("PostSocket Handle (CreateComment) error: %w", err)
					}
				}
			}
		}
	return nil
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
		Type:  post,
		Posts: p,
	}

	return c.Broadcast(s)
}

func CreatePost(post *database.Post) error {
	stmt, err := database.DB.Prepare("INSERT INTO posts (postID, nickname, title, categories, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	if post.PostID == "" {
		post.PostID = uuid.NewV4().String()
	}

	// TODO: remove placeholder nickname once login/sessions are working
	if post.Nickname == "" {
		post.Nickname = "Cassidy"
	}

	_, err = stmt.Exec(post.PostID, post.Nickname, post.Title, post.Categories, post.Body)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return nil
}

func CreateComment(comment database.Comment) error {
	stmt, err := database.DB.Prepare("INSERT INTO comments (commentID, postID, nickname, commentText) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	if comment.CommentID == "" {
		comment.CommentID = uuid.NewV4().String()
	}

	// TODO: remove placeholder username once login/sessions are working
	if comment.Nickname == "" {
		comment.Nickname = "Cassidy"
	}

	_, err = stmt.Exec(comment.CommentID, comment.PostID, comment.Nickname, comment.Body)
	if err != nil {
		return fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	return nil
}
