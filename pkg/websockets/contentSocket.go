package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"strings"
	"text/template"
	"time"
)

type ContentMessage struct {
	Type      messageType `json:"type,omitempty"`
	Body      string      `json:"body,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Nickname  string      `json:"nickname,omitempty"`
	Resource  string      `json:"resource,omitempty"`
	PostID    string      `json:"post_id,omitempty"`
}

func (m *ContentMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (content) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send content message down ws of type %s", s.t.String())
	}
	return nil
}

func OnContentConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	tpl, err := template.ParseGlob("templates/*")
	if err != nil {
		return err
	}

	sb := &strings.Builder{}

	if err := tpl.ExecuteTemplate(sb, "home.template", nil); err != nil {
		return fmt.Errorf("Home ExecuteTemplate error: %w", err)
	}

	c := &ContentMessage{
		Type: content,
		Body: sb.String(),
	}

	return c.Broadcast(s)
}

func (m *ContentMessage) Handle(s *socket) error {
	tpl, err := template.ParseGlob("templates/*")
	if err != nil {
		return err
	}

	sb := &strings.Builder{}
	switch m.Resource {
	case "post":
		if err := tpl.ExecuteTemplate(sb, "home.template", nil); err != nil {
			return fmt.Errorf("Home ExecuteTemplate error: %w", err)
		}
	case "profile":
		if err := tpl.ExecuteTemplate(sb, "profile.template", nil); err != nil {
			return fmt.Errorf("Profile ExecuteTemplate error: %w", err)
		}
	case "login":
		if err := tpl.ExecuteTemplate(sb, "reg-log.template", nil); err != nil {
			return fmt.Errorf("Reg-Log ExecuteTemplate error: %+v\n", err)
		}
	case "presence":
		if err := tpl.ExecuteTemplate(sb, "presence.template", nil); err != nil {
			return fmt.Errorf("Presence ExecuteTemplate error: %+v\n", err)
		}
	case "comment":
		if m.PostID == "" {
			return fmt.Errorf("Empty post ID when requesting comments")
		}
		comments, err := database.GetComments()
		if err != nil {
			return fmt.Errorf("Unable to get comments for comment template: %w", err)
		}
		comments = database.FilterCommentsForPost(m.PostID, comments)
		allPosts, err := database.GetPosts()
		if err != nil {
			return fmt.Errorf("Unable to get posts for comment template: %w", err)
		}
		newPost := database.Post{}
		for _, pst := range allPosts {
			if pst.PostID == m.PostID {
				newPost = *pst
			}
		}
		// post, err := database.GetPostForComment(comments[0])
		// if err != nil {
		// 	return fmt.Errorf("Unable to get post for comments: %w", err)
		// }

		commentTemplateData := struct {
			Post     database.Post
			Comments []database.Comment
		}{
			Post:     newPost,
			Comments: comments,
		}

		if err := tpl.ExecuteTemplate(sb, "comment.template", commentTemplateData); err != nil {
			return fmt.Errorf("Comment ExecuteTemplate error: %+v\n", err)
		}
	default:
		return fmt.Errorf("template %s not found", m.Resource)
	}
	m.Body = sb.String()
	return m.Broadcast(s)
}
