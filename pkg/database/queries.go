package database

import (
	"fmt"
)

type User struct {
	ID        string `json:"id,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Age       string `json:"age,omitempty"`
	Gender    string `json:"gender,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type Post struct {
	PostID     string    `json:"post_id,omitempty"`
	Nickname   string    `json:"nickname"`
	Title      string    `json:"title,omitempty"`
	Categories string    `json:"categories,omitempty"`
	Body       string    `json:"body,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
	Nickname  string `json:"nickname"`
	Body      string `json:"body,omitempty"`
}

type Conversation struct {
	ChatID       string   `json:"chat_id,omitempty"`
	Participants []string `json:"participants"`
	Chats        []Chat   `json:"chats,omitempty"`
}

type Chat struct {
	ChatID   string `json:"chat_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Date     string `json:"date,omitempty"`
	Body     string `json:"body,omitempty"`
}

func GetUsers() ([]User, error) {
	users := []User{}
	rows, err := DB.Query(`SELECT * FROM users`)
	if err != nil {
		return users, fmt.Errorf("GetUsers DB Query error: %+v\n", err)
	}
	var id string
	var nickname string
	var age string
	var gender string
	var firstname string
	var lastname string
	var email string
	var password string

	for rows.Next() {
		err := rows.Scan(&id, &nickname, &age, &gender, &firstname, &lastname, &email, &password)
		if err != nil {
			return users, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		users = append(users, User{
			ID:        id,
			Nickname:  nickname,
			Age:       age,
			Gender:    gender,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Password:  password,
		})
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}

func GetPosts() ([]*Post, error) {
	posts := []*Post{}
	rows, err := DB.Query(`SELECT * FROM posts`)
	if err != nil {
		return posts, fmt.Errorf("GetPosts DB Query error: %+v\n", err)
	}
	var postid string
	var nickname string
	var title string
	var category string
	var postcontent string

	for rows.Next() {
		err := rows.Scan(&postid, &nickname, &title, &category, &postcontent)
		if err != nil {
			return posts, fmt.Errorf("GetPosts rows.Scan error: %+v\n", err)
		}
		posts = append(posts, &Post{
			PostID:     postid,
			Nickname:   nickname,
			Title:      title,
			Categories: category,
			Body:       postcontent,
		})
	}
	err = rows.Err()
	if err != nil {
		return posts, err
	}
	return posts, nil
}

func GetPostForComment(c Comment) (Post, error) {
	posts, err := GetPosts()
	if err != nil {
		return Post{}, err
	}
	for _, p := range posts {
		if p.PostID == c.PostID {
			return *p, nil
		}
	}
	return Post{}, fmt.Errorf("no matching post found for id: %s", c.PostID)
}

func GetPopulatedPosts() ([]*Post, error) {
	posts, err := GetPosts()
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedPosts (GetPosts) error: %+v\n", err)
	}

	populatedPosts, err := populateCommentsForPosts(posts)
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedPosts (populateCommentsForPosts) error: %+v\n", err)
	}

	return populatedPosts, nil
}

func populateCommentsForPosts(posts []*Post) ([]*Post, error) {
	comments, err := GetComments()
	if err != nil {
		return nil, fmt.Errorf("populateCommentsForPosts (GetComments) error: %+v\n", err)
	}
	outPost := []*Post{}
	for _, pts := range posts {
		newPost := pts
		for _, cmts := range comments {
			if pts.PostID == cmts.PostID {
				newPost.Comments = append(newPost.Comments, cmts)
			}
		}
		outPost = append(outPost, newPost)
	}
	return outPost, nil
}

func GetComments() ([]Comment, error) {
	comments := []Comment{}
	rows, err := DB.Query(`SELECT * FROM comments`)
	if err != nil {
		return comments, fmt.Errorf("GetComments DB Query error: %+v\n", err)
	}
	var postid string
	var commentid string
	var nickname string
	var commentcontent string

	for rows.Next() {
		err := rows.Scan(&commentid, &postid, &nickname, &commentcontent)
		if err != nil {
			return comments, fmt.Errorf("GetComments rows.Scan error: %+v\n", err)
		}
		comments = append(comments, Comment{
			CommentID: commentid,
			PostID:    postid,
			Nickname:  nickname,
			Body:      commentcontent,
		})
	}
	err = rows.Err()
	if err != nil {
		return comments, err
	}
	return comments, nil
}

func FilterCommentsForPost(postID string, comments []Comment) []Comment {
	out := []Comment{}
	for _, c := range comments {
		if postID == c.PostID {
			out = append(out, c)
		}
	}
	return out
}

func GetConversations() ([]*Conversation, error) {
	conversations := []*Conversation{}
	rows, err := DB.Query(`SELECT * FROM conversations`)
	if err != nil {
		return conversations, fmt.Errorf("GetConversations DB Query error: %+v\n", err)
	}
	var chatid string
	participants := []string{}

	for rows.Next() {
		err := rows.Scan(&chatid, &participants)
		if err != nil {
			return conversations, fmt.Errorf("GetConversations rows.Scan error: %+v\n", err)
		}
		conversations = append(conversations, &Conversation{
			ChatID:       chatid,
			Participants: participants,
		})
	}
	err = rows.Err()
	if err != nil {
		return conversations, err
	}
	return conversations, nil
}

func GetChat() ([]*Chat, error) {
	chat := []*Chat{}
	rows, err := DB.Query(`SELECT * FROM chat`)
	if err != nil {
		return chat, fmt.Errorf("GetChat DB Query error: %+v\n", err)
	}
	var chatid string
	var sender string
	var receiver string
	var date string
	var body string

	for rows.Next() {
		err := rows.Scan(&chatid, &sender, &receiver, &date, &body)
		if err != nil {
			return chat, fmt.Errorf("GetChat rows.Scan error: %+v\n", err)
		}
		chat = append(chat, &Chat{
			ChatID:   chatid,
			Sender:   sender,
			Receiver: receiver,
			Date:     date,
			Body:     body,
		})
	}
	err = rows.Err()
	if err != nil {
		return chat, err
	}
	return chat, nil
}

func populateChatForConversation(chats []*Conversation) ([]*Conversation, error) {
	chat, err := GetChat()
	if err != nil {
		return nil, fmt.Errorf("populateChatForChats (GetChat) error: %+v\n", err)
	}
	outChat := []*Conversation{}
	outParticipants := []string{}
	for _, cts := range chats {
		newChat := cts
		for _, cht := range chat {
			if cts.ChatID == cht.ChatID {
				outParticipants = append(outParticipants, cht.Receiver)
				outParticipants = append(outParticipants, cht.Sender)
				newChat.Participants = append(newChat.Participants, outParticipants...)
			}
		}
		outChat = append(outChat, newChat)
	}
	return outChat, nil
}

func GetPopulatedConversations() ([]*Conversation, error) {
	conversations, err := GetConversations()
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedConversation (GetConversations) error: %+v\n", err)
	}

	populatedConversations, err := populateChatForConversation(conversations)
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedConversation (populateChatForConversation) error: %+v\n", err)
	}

	return populatedConversations, nil
}
