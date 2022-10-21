package database

import "fmt"

type User struct {
	ID       string
	Email    string
	Username string
	Password string
	Nickname string
}

type Post struct {
	PostID     string    `json:"post_id,omitempty"`
	Username   string    `json:"username,omitempty"`
	Title      string    `json:"text,omitempty"`
	Categories string    `json:"categories,omitempty"`
	Body       string    `json:"body,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Body      string `json:"body,omitempty"`
}

func GetUsers() ([]User, error) {
	users := []User{}
	rows, err := DB.Query(`SELECT * FROM users`)
	if err != nil {
		return users, fmt.Errorf("GetUsers DB Query error: %+v\n", err)
	}
	var id string
	var email string
	var username string
	var password string
	var nickname string
	for rows.Next() {
		err := rows.Scan(&id, &email, &username, &password, &nickname)
		if err != nil {
			return users, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		users = append(users, User{
			ID:       id,
			Email:    email,
			Username: username,
			Password: password,
			Nickname: nickname,
		})
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}

func GetPosts() ([]Post, error) {
	posts := []Post{}
	rows, err := DB.Query(`SELECT * FROM posts`)
	if err != nil {
		return posts, fmt.Errorf("GetPosts DB Query error: %+v\n", err)
	}
	var postid string
	var username string
	var category string
	var title string
	var postcontent string
	for rows.Next() {
		err := rows.Scan(&postid, &username, &category, &title, &postcontent)
		if err != nil {
			return posts, fmt.Errorf("GetPosts rows.Scan error: %+v\n", err)
		}
		posts = append(posts, Post{
			PostID:     postid,
			Username:   username,
			Categories: category,
			Title:      title,
			Body:       postcontent,
		})
	}
	err = rows.Err()
	if err != nil {
		return posts, err
	}
	return posts, nil
}
