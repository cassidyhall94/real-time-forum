package database

import "fmt"

type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Age      string `json:"age,omitempty"`
}

type Post struct {
	PostID     string `json:"post_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Title      string `json:"title,omitempty"`
	Categories string `json:"categories,omitempty"`
	Body       string `json:"body,omitempty"`
	// Comments   []Comment `json:"comments,omitempty"`
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
	var gender string
	var age string

	for rows.Next() {
		err := rows.Scan(&id, &email, &username, &password, &nickname, &gender, &age)
		if err != nil {
			return users, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		users = append(users, User{
			ID:       id,
			Email:    email,
			Username: username,
			Password: password,
			Nickname: nickname,
			Gender:   gender,
			Age:      age,
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
