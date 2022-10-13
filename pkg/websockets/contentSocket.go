package websockets

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

type ContentMessage struct {
	Type      messageType `json:"type,omitempty"`
	Body      string      `json:"body,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Username  string      `json:"username,omitempty"`
	Resource  string      `json:"resource,omitempty"`
}

func (m *ContentMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send content message down ws of type %s", s.t.String())
	}
	return nil
}

func OnContentConnect(s *socket) error {
	tpl, err := template.ParseGlob("templates/*")
	if err != nil {
		return err
	}

	w, err := s.con.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if err := tpl.ExecuteTemplate(w, "home.template", nil); err != nil {
		return fmt.Errorf("Home ExecuteTemplate error: %w", err)
	}

	return w.Close()
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
	default:
		return fmt.Errorf("template %s not found", m.Resource)
	}
	m.Body = sb.String()
	return m.Broadcast(s)
}
