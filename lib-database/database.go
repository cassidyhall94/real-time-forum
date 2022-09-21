package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver connects go with sql
)

var DB *sql.DB

func (forum *Forum) CreateUser(user User) error {
	stmt, err := forum.DB.Prepare("INSERT INTO people (uuid, username, email, password) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateUser DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(user.Uuid, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("CreateUser Exec error: %+v\n", err)
	}
	return nil
}

func (forum *Forum) CreateSession(session Session) error {
	stmt, err := forum.DB.Prepare("INSERT INTO session (sessionID, userName, expiryTime) VALUES (?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateSession DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(session.SessionID, session.Username, session.Expiry)
	if err != nil {
		return fmt.Errorf("CreateSession Exec error: %+v\n", err)
	}
	return nil
}

func (forum *Forum) CreatePost(post PostFeed) error {
	stmt, err := forum.DB.Prepare("INSERT INTO post (username, title, content, category, dateCreated) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(post.Username, post.Title, post.Content, post.Category, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return nil
}

func (forum *Forum) CreateReaction(reaction Reaction) error {
	stmt, err := forum.DB.Prepare("INSERT INTO reaction (postid, username, commentid, liked, disliked) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateReaction DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(reaction.PostID, reaction.Username, reaction.CommentID, reaction.Liked, reaction.Disliked)
	if err != nil {
		return fmt.Errorf("CreateReactions Exec error: %+v\n", err)
	}
	return nil
}

func (forum *Forum) CreateComment(comment Comment) error {
	stmt, err := forum.DB.Prepare("INSERT INTO comments ( postID, userID, content, dateCreated) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec(comment.PostID, comment.UserId, comment.Content, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	return nil
}

// Update(Updates an item in a table).
func (feed *Forum) UpdatePost(item PostFeed) error {
	stmt, err := feed.DB.Prepare("UPDATE post SET title = ?, content = ?, category = ? WHERE postID = ?;")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("UpdatePost DB Prepare error: %+v", err)
	}
	// stmt.QueryRow(stmt, item.Title, item.Content, item.Category)
	_, err = stmt.Exec(item.Title, item.Content, item.Category, item.PostID)
	if err != nil {
		return fmt.Errorf("unable to insert item into post: %w", err)
	}
	return nil
}

func (feed *Forum) UpdateReaction(item Reaction) error {
	stmt, err := feed.DB.Prepare("UPDATE reaction SET liked = ?, disliked = ? WHERE (postid = ? AND username = ? AND reactionID = ?) OR (commentid = ? AND username = ? AND reactionID = ?);")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("UpdateReaction DB Prepare error: %+v", err)
	}
	_, err = stmt.Exec(item.Liked, item.Disliked, item.PostID, item.Username, item.ReactionID, item.CommentID, item.Username, item.ReactionID)
	if err != nil {
		return fmt.Errorf("unable to update reaction: %w", err)
	}
	return nil
}

func (feed *Forum) AssertUniqueSessionForUser(session Session) error {
	stmt, err := feed.DB.Prepare("DELETE FROM session WHERE userName == ? AND sessionID != ?;")
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("AssertUniqueSessionForUser DB Prepare error: %w", err)
	}
	_, err = stmt.Exec(session.Username, session.SessionID)
	if err != nil {
		return fmt.Errorf("unable to assert unique session: %w", err)
	}
	_, err = feed.Exec("PRAGMA wal_checkpoint(FULL);")
	if err != nil {
		return fmt.Errorf("unable to confirm full write of unique session to db: %w", err)
	}
	return nil
}

// ---------------------------------------------- TABLES ------------------------------- --//

func userTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS people (
	userID INTEGER PRIMARY KEY,	
	uuid TEXT, 
	username TEXT UNIQUE,
	email TEXT UNIQUE, 
	password TEXT);
`)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("userTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("userTable Exec error: %+v\n", err)
	}
	return nil
}

func sessionTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS session (
	sessionID TEXT PRIMARY KEY REFERENCES people(uuid),	
	userName TEXT REFERENCES people(username),
	expiryTime TEXT);
	`)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("sessionTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("sessionTable Exec error: %+v\n", err)
	}
	return nil
}

func postTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS post (
 postID INTEGER PRIMARY KEY AUTOINCREMENT,
 username TEXT REFERENCES session(userName),
 title TEXT,
 content TEXT, 
 category TEXT,
 dateCreated TEXT);
 `)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("postTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("postTable Exec error: %+v\n", err)
	}
	return nil
}

// @TODO: add likes/dislikes to comments.
func commentTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS comments (
   commentID INTEGER PRIMARY KEY AUTOINCREMENT, 
   postID INTEGER REFERENCES people(userID), 
	userID STRING REFERENCES session(userName),
	content TEXT NOT NULL, 
	dateCreated TEXT NOT NULL);
	`)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("commentTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("commentTable Exec error: %+v\n", err)
	}
	return nil
}

func reactionsTable(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS reaction (
   reactionID INTEGER PRIMARY KEY AUTOINCREMENT,
   postID INTEGER REFERENCES posts(postID),
   username TEXT REFERENCES people(username),
   commentID INTEGER REFERENCES comments(commentID),
   liked BOOL,
   disliked BOOL);
	`)
	defer stmt.Close()
	if err != nil {
		return fmt.Errorf("reactionsTable DB Prepare error: %+v\n", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("reactionsTable Exec error: %+v\n", err)
	}
	return nil
}

func Connect(db *sql.DB) (*Forum, error) {
	err := userTable(db)
	if err != nil {
		return nil, fmt.Errorf("Connect userTable error: %+v\n", err)
	}
	err2 := sessionTable(db)
	if err != nil {
		return nil, fmt.Errorf("Connect sessionTable error: %+v\n", err2)
	}
	err3 := postTable(db)
	if err != nil {
		return nil, fmt.Errorf("Connect postTable error: %+v\n", err3)
	}
	err4 := commentTable(db)
	if err != nil {
		return nil, fmt.Errorf("Connect commentTable error: %+v\n", err4)
	}
	err5 := reactionsTable(db)
	if err != nil {
		return nil, fmt.Errorf("Connect reactionsTable error: %+v\n", err5)
	}

	return &Forum{
		DB: db,
	}, nil
}

func (data *Forum) GetPosts() ([]PostFeed, error) {
	posts := []PostFeed{}
	rows, err := data.DB.Query(`SELECT * FROM post`)
	if err != nil {
		return posts, fmt.Errorf("GetPosts DB Query error: %+v\n", err)
	}

	var id int
	var uiD string
	var title string
	var content string
	var created string
	var category string

	for rows.Next() {
		err := rows.Scan(&id, &uiD, &title, &content, &category, &created)
		if err != nil {
			return posts, fmt.Errorf("GetPosts rows.Scan error: %+v\n", err)
		}

		likes, err := getLikesForPost(data.DB, id)
		if err != nil {
			return posts, fmt.Errorf("GetPosts getLikesForPost error: %+v\n", err)
		}

		dislikes, err := getDislikesForPost(data.DB, id)
		if err != nil {
			return posts, fmt.Errorf("GetPosts getDislikesForPost error: %+v\n", err)
		}

		posts = append(posts, PostFeed{
			PostID:    id,
			Username:  uiD,
			Title:     title,
			Content:   content,
			Category:  category,
			CreatedAt: created,
			Likes:     likes,
			Dislikes:  dislikes,
		})
	}

	// fmt.Println(posts)
	return posts, nil
}

func getLikesForPost(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT liked FROM reaction WHERE liked = TRUE AND postID = ?")
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("getLikesForPost DB Prepare error: %+v\n", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getLikesForPost DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

func getDislikesForPost(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT disliked FROM reaction WHERE disliked = TRUE AND postID = ?")
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("getDislikesForPost DB Prepare error: %+v\n", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getDislikesForPost DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

func getLikesForComment(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT liked FROM reaction WHERE liked = TRUE AND commentID = ?")
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("getLikesForComment DB Prepare error: %+v\n", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getLikesForComment DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

func getDislikesForComment(db *sql.DB, id int) (int, error) {
	stmt, err := db.Prepare("SELECT disliked FROM reaction WHERE disliked = TRUE AND commentID = ?")
	defer stmt.Close()
	if err != nil {
		return 0, fmt.Errorf("getDislikesForComment DB Prepare error: %+v\n", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return 0, fmt.Errorf("getDislikesForComment DB Query error: %+v\n", err)
	}

	counter := 0
	for rows.Next() {
		counter++
	}
	return counter, nil
}

func (data *Forum) GetReactions() ([]Reaction, error) {
	reactions := []Reaction{}
	rows, err := data.DB.Query(`SELECT * FROM reaction`)
	if err != nil {
		return reactions, fmt.Errorf("GetReactions DB Query error: %+v\n", err)
	}

	var reactionID int
	var postID int
	var username string
	var commentID int
	var liked bool
	var disliked bool

	for rows.Next() {
		err := rows.Scan(&reactionID, &postID, &username, &commentID, &liked, &disliked)
		if err != nil {
			return reactions, fmt.Errorf("GetReactions rows.Scan error: %+v\n", err)
		}
		reactions = append(reactions, Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			Username:   username,
			CommentID:  commentID,
			Liked:      liked,
			Disliked:   disliked,
		})
	}
	// fmt.Println(reactions)
	return reactions, nil
}

func (data *Forum) GetReactionByPostID(targetPostID, targetUsername string) (*Reaction, error) {
	stmt, err := data.DB.Prepare("SELECT * FROM reaction WHERE postID = ? AND username = ?")
	defer stmt.Close()
	if err != nil {
		return nil, fmt.Errorf("GetReactionByPostID DB Prepare error: %+v", err)
	}
	rows, err := stmt.Query(targetPostID, targetUsername)
	if err != nil {
		return nil, fmt.Errorf("GetReactionByPostID DB Query error: %+v", err)
	}

	var reactionID int
	var postID int
	var username string
	var commentID int
	var liked bool
	var disliked bool

	for rows.Next() {
		err := rows.Scan(&reactionID, &postID, &username, &commentID, &liked, &disliked)
		if err != nil {
			return nil, fmt.Errorf("GetReactionByPostID rows.Scan error: %+v\n", err)
		}
		return &Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			Username:   username,
			CommentID:  commentID,
			Liked:      liked,
			Disliked:   disliked,
		}, nil
	}
	return nil, nil
}

func (data *Forum) GetReactionByCommentID(targetCommentID, targetUsername string) (*Reaction, error) {
	stmt, err := data.DB.Prepare("SELECT * FROM reaction WHERE commentID = ? AND username = ?")
	defer stmt.Close()
	if err != nil {
		return nil, fmt.Errorf("GetReactionByCommentID DB Prepare error: %+v", err)
	}
	rows, err := stmt.Query(targetCommentID, targetUsername)
	if err != nil {
		return nil, fmt.Errorf("GetReactionByCommentID DB Query error: %+v", err)
	}

	var reactionID int
	var postID int
	var username string
	var commentID int
	var liked bool
	var disliked bool

	for rows.Next() {
		err := rows.Scan(&reactionID, &postID, &username, &commentID, &liked, &disliked)
		if err != nil {
			return nil, fmt.Errorf("GetReactionByCommentID rows.Scan error: %+v\n", err)
		}
		return &Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			Username:   username,
			CommentID:  commentID,
			Liked:      liked,
			Disliked:   disliked,
		}, nil
	}
	return nil, nil
}

// @TODO: add likes/dislikes(reactions) to comments.
func (data *Forum) GetComments() ([]Comment, error) {
	comments := []Comment{}
	rows, err := data.DB.Query(`SELECT * FROM comments`)
	if err != nil {
		return comments, fmt.Errorf("GetComments DB Query error: %+v\n", err)
	}
	var commentid int
	var postid int
	var userid string
	var content string
	var created string

	for rows.Next() {
		err := rows.Scan(&commentid, &postid, &userid, &content, &created)
		if err != nil {
			return comments, fmt.Errorf("GetComments rows.Scan error: %+v\n", err)
		}
		likes, err := getLikesForComment(data.DB, commentid)
		if err != nil {
			return comments, fmt.Errorf("GetComments getLikesForComment error: %+v\n", err)
		}

		dislikes, err := getDislikesForComment(data.DB, commentid)
		if err != nil {
			return comments, fmt.Errorf("GetComments getDislikesForComment error: %+v\n", err)
		}

		comments = append(comments, Comment{
			CommentID: commentid,
			PostID:    postid,
			UserId:    userid,
			Content:   content,
			CreatedAt: created,
			Likes:     likes,
			Dislikes:  dislikes,
		})
	}
	return comments, nil
}

func (data *Forum) GetSessions() ([]Session, error) {
	session := []Session{}
	rows, err := data.DB.Query(`SELECT * FROM session`)
	if err != nil {
		return session, fmt.Errorf("GetSession DB Query error: %+v\n", err)
	}
	var session_token string
	var uName string
	var exTime string

	for rows.Next() {
		err := rows.Scan(&session_token, &uName, &exTime)
		if err != nil {
			return nil, fmt.Errorf("GetSession rows.Scan error: %+v\n", err)
		}
		// time.Format("01-02-2006 15:04")
		convTime, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", exTime)
		if err != nil {
			return nil, fmt.Errorf("GetSession time.Parse(exTime) error: %+v\n", err)
		}
		session = append(session, Session{
			SessionID: session_token,
			Username:  uName,
			Expiry:    convTime,
		})
	}
	return session, nil
}
