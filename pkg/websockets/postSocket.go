package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PostMessage struct {
	Type      messageType     `json:"type,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
	Return    string          `json:"return,omitempty"`
	Posts     []database.Post `json:"posts,omitempty"`
}

// TODO: add code for handling comments and attaching to post
func (m PostMessage) Handle(s *socket) error {
	fmt.Println(m.Return)
	if m.Return == "all posts" {
		if err := OnPostsConnect(s); err != nil {
			return err
		}
		m.Return = ""
	} else {
		for _, post := range m.Posts {
			if m.Return == post.PostID {
				if err := OnCommentsConnect(s); err != nil {
					return err
				}
				m.Return = ""
			} else {
				if err := CreatePost(post); err != nil {
					return err
				}
			}
		}
		return m.Broadcast(s)
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

// TODO: add timestamp
func OnCommentsConnect(s *socket) error {
	posts, err := database.GetPosts()
	if err != nil {
		return fmt.Errorf("OnPostsConnect (GetPosts) error: %+v\n", err)
	}

	for _, pts := range posts {
		comments, err := database.GetComments()
		if err != nil {
			return fmt.Errorf("OnCommentsConnect (GetComments) error: %+v\n", err)
		}
		c := &PostMessage{
			Type:      post,
			Timestamp: "",
			Return:    "all posts",
			Posts: []database.Post{
				{
					PostID:     pts.PostID,
					Username:   pts.Username,
					Title:      pts.Title,
					Categories: pts.Categories,
					Body:       pts.Body,
					Comments:   comments,
				},
			},
		}
		c.Broadcast(s)
	}
	return nil
}

// TODO: add timestamp
func OnPostsConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	posts, err := database.GetPosts()
	if err != nil {
		return fmt.Errorf("OnPostsConnect (GetPosts) error: %+v\n", err)
	}
	c := &PostMessage{
		Type:      post,
		Timestamp: "",
		Return:    "all posts",
		Posts:     posts,
	}
	return c.Broadcast(s)
}

func CreatePost(post database.Post) error {
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
