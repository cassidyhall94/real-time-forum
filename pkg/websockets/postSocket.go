package websockets

import (
	"database/sql"
	"fmt"
	"log"
	auth "real-time-forum/pkg/authentication"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type PostMessage struct {
	Type       messageType `json:"type,omitempty"`
	Header     string      `json:"text,omitempty"`
	Body       string      `json:"body,omitempty"`
	Categories []string    `json:"categories,omitempty"`
	Timestamp  string      `json:"timestamp,omitempty"`
	Username   string      `json:"username,omitempty"`
}

func (m PostMessage) Handle(s *socket) error {
	return nil
}

func (m *PostMessage) Broadcast() error {
	for _, s := range savedSockets {
		if s.t == m.Type {
			if err := s.con.WriteJSON(m); err != nil {
				return fmt.Errorf("unable to send (post )message: %w", err)
			}
		}
	}
	return nil
}

type postDisplay struct {
	PostID        string          `json:"post_id,omitempty"`
	Username      string          `json:"username,omitempty"`
	PostCategory  string          `json:"post_category,omitempty"`
	TitleText     string          `json:"title_text,omitempty"`
	PostText      string          `json:"post_text,omitempty"`
	CookieChecker bool            `json:"cookie_checker,omitempty"`
	Comments      []commentStruct `json:"comments,omitempty"`
}

type commentStruct struct {
	CommentID       string `json:"comment_id,omitempty"`
	CpostID         string `json:"cpost_id,omitempty"`
	CommentUsername string `json:"comment_username,omitempty"`
	CommentText     string `json:"comment_text,omitempty"`
	CookieChecker   bool   `json:"cookie_checker,omitempty"`
}

func newPost(userName, category, title, post string, db *sql.DB) {
	if title == "" {
		return
	}
	fmt.Println("ADDING POST")
	uuid := uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO posts (postID, userName, category, title, post) VALUES (?, ?, ?,  ?, ?)", uuid, userName, category, title, post)
	if err != nil {
		fmt.Println("Error adding new post")
		log.Fatal(err.Error())
	}
	auth.Person.PostAdded = true
	catSlc := strings.Split(category, " ")
	golangSelected := 0
	javascriptSelected := 0
	rustSelected := 0
	for _, r := range catSlc {
		if r == "Golang" {
			golangSelected = 1
		} else if r == "Javascript" {
			javascriptSelected = 1
		} else if r == "Rust" {
			rustSelected = 1
		}
	}
	_, errAddCats := db.Exec("INSERT INTO categories (postID, Golang, Javascript, Rust) VALUES (?, ?, ?, ?)", uuid, golangSelected, javascriptSelected, rustSelected)
	if errAddCats != nil {
		fmt.Println("ERROR when adding into the category table")
	}
}

func postData(db *sql.DB) []postDisplay {
	rows, err := db.Query("SELECT postID, userName, category, title, post FROM posts")
	if err != nil {
		fmt.Println("Error selecting post data")
		log.Fatal(err.Error())
	}
	finalArray := []postDisplay{}
	for rows.Next() {
		var u postDisplay
		err := rows.Scan(
			&u.PostID,
			&u.Username,
			&u.PostCategory,
			&u.TitleText,
			&u.PostText,
		)
		u.CookieChecker = auth.Person.CookieChecker
		if err != nil {
			fmt.Println("SCANNING ERROR")
			log.Fatal(err.Error())
		}
		commentSlc := []commentStruct{}
		var tempComStruct commentStruct
		commentRow, errComs := db.Query("SELECT commentID, postID, username, commentText  FROM comments WHERE postID = ?", u.PostID)
		if errComs != nil {
			fmt.Println("Error selecting comment data")
			log.Fatal(errComs.Error())
		}
		for commentRow.Next() {
			err := commentRow.Scan(
				&tempComStruct.CommentID,
				&tempComStruct.CpostID,
				&tempComStruct.CommentUsername,
				&tempComStruct.CommentText,
			)
			tempComStruct.CookieChecker = auth.Person.CookieChecker
			if err != nil {
				fmt.Println("Error scanning comments")
				log.Fatal(err.Error())
			}
			fmt.Printf("\nCOMMENT STRUCT_____-------------------------------------%v\n\n", tempComStruct)
			commentSlc = append(commentSlc, tempComStruct)
		}
		u.Comments = commentSlc
		finalArray = append(finalArray, u)
		for i, j := 0, len(finalArray)-1; i < j; i, j = i+1, j-1 {
			finalArray[i], finalArray[j] = finalArray[j], finalArray[i]
		}
	}
	return finalArray
}

func newComment(userName, postID, commentText string, db *sql.DB) {
	if commentText == "" {
		return
	}
	fmt.Println("ADDING Comment")
	uuid := uuid.NewV4().String()
	_, err := db.Exec("INSERT INTO comments (commentID, postID, userName, commentText) VALUES (?, ?, ?, ?)", uuid, postID, userName, commentText)
	if err != nil {
		fmt.Println("ERROR ADDING COMMENT TO THE TABLE")
		log.Fatal(err.Error())
	}
	auth.Person.PostAdded = true
}

func PostGetter(postIDSlc []string, db *sql.DB) []postDisplay {
	finalArray := []postDisplay{}
	for _, r := range postIDSlc {
		rows, errDetails := db.Query("SELECT postID, userName, category, title, post FROM posts WHERE postID = (?)", r)
		if errDetails != nil {
			fmt.Println("ERROR when selecting the information for specific posts (func POSTGETTER)")
			log.Fatal(errDetails.Error())
		}
		for rows.Next() {
			var postDetails postDisplay
			err := rows.Scan(
				&postDetails.PostID,
				&postDetails.Username,
				&postDetails.PostCategory,
				&postDetails.TitleText,
				&postDetails.PostText,
			)
			postDetails.CookieChecker = auth.Person.CookieChecker
			if err != nil {
				fmt.Println("ERROR Scanning through the rows (func POSTGETTER)")
				log.Fatal(err.Error())
			}
			finalArray = append(finalArray, postDetails)
		}
	}
	return finalArray
}
