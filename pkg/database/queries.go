package database

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
	_ "golang.org/x/exp/constraints"
)

type User struct {
	ID        string `json:"id,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Age       string `json:"age,omitempty"`
	Gender    string `json:"gender,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"-"`
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
	ConvoID      string `json:"convo_id"`
	Participants []User `json:"participants"`
	Chats        []Chat `json:"chats"`
}

type Chat struct {
	ConvoID string `json:"convo_id"`
	ChatID  string `json:"chat_id"`
	Sender  User   `json:"sender"`
	Date    string `json:"date,omitempty"`
	Body    string `json:"body"`
}

type Presence struct {
	User              *User     `json:"user,omitempty"`
	Online            bool      `json:"online,omitempty"`
	LastContactedTime time.Time `json:"last_contacted_time,omitempty"`
}

type Login struct {
	Nickname string `json:"nickname,omitempty"`
	Password string `json:"password,omitempty"`
}

type Session struct {
	SessionID  string
	UserID     string
	ExpiryTime string
}

// type Cookie struct {
// 	Name string
// 	Value string
// 	Path string
// 	Domain string
// 	Expires time.Time
// 	RawExpires string
// 	MaxAge int
// 	Secure bool
// 	HttpOnly bool
// 	Raw string
// 	Unparsed []string
// }

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

func GetSessions() ([]*Session, error) {
	sessions := []*Session{}

	rows, err := DB.Query(` SELECT * FROM sessions`)
	if err != nil {
		return sessions, fmt.Errorf("GetSessions DB Query error: %+v\n", err)
	}

	var sessionID string
	var userID string
	var expiryTime string

	for rows.Next() {
		err := rows.Scan(&sessionID, &userID, &expiryTime)
		if err != nil {
			return sessions, fmt.Errorf("GetSessions rows.Scan error: %+v\n", err)
		}
		sessions = append(sessions, &Session{
			SessionID:  sessionID,
			UserID:     userID,
			ExpiryTime: expiryTime,
		})
	}
	err = rows.Err()
	if err != nil {
		return sessions, err
	}
	return sessions, nil
}

// GetPresences marries users and sessions to create map of users split by their being logged in or out of the forum
func GetPresencesForUser(user *User) ([]*Presence, error) {
	users, err := GetUsers()
	if err != nil {
		return nil, fmt.Errorf("unable to get users for presenses: %w", err)
	}
	sessions, err := GetSessions()
	if err != nil {
		return nil, fmt.Errorf("unable to get sessions for presenses: %w", err)
	}

	out := []*Presence{}
	for _, u := range users {
		lc, err := GetLastContactedForUsers(user.ID, u.ID)
		if err != nil {
			return nil, fmt.Errorf("error looking for last contacted between '%s' and '%s': %w", user.ID, u.ID, err)
		}
		if lc == nil {
			lc = &time.Time{}
		}
		k := u
		p := &Presence{
			User:              &k,
			Online:            false,
			LastContactedTime: *lc,
		}
		for _, s := range sessions {
			if s.UserID == u.ID {
				p.Online = true
			}
		}
		out = append(out, p)
	}

	return out, nil
}

func GetUserFromSessionID(sess *Session) (*User, error) {
	users, err := GetUsers()
	if err != nil {
		return nil, fmt.Errorf("unable to get users for GetUserFromSessionID: %w", err)
	}
	sessions, err := GetSessions()
	if err != nil {
		return nil, fmt.Errorf("unable to get sessions for GetUserFromSessionID: %w", err)
	}

	for _, u := range users {
		for _, s := range sessions {
			if u.ID == s.UserID {
				return &u, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to find user from sessionID '%s'", sess.SessionID)
}

func GetLastContactedForUsers(userIDs ...string) (*time.Time, error) {
	convos, err := GetConversations()
	if err != nil {
		return nil, fmt.Errorf("unable to get conversations: %w", err)
	}
	cid, err := GetConvoID(userIDs, convos)
	if err != nil {
		return nil, fmt.Errorf("unable to get convoid: %w", err)
	}

	var convo *Conversation = nil
	for _, c := range convos {
		if c.ConvoID == cid {
			convo = c
			if err := populateChatsForConversation(c); err != nil {
				return nil, fmt.Errorf("unable to populate chats in convo '%s': %w", cid, err)
			}
		}
	}
	if convo == nil {
		return nil, nil
	}


	var d *time.Time = &time.Time{}
	for i, ch := range convo.Chats {
		t, err := time.Parse("2006-01-02 15:04:05-07:00", ch.Date)
		if err != nil {
			return d, fmt.Errorf("unable to parse time for chat '%s': %w", ch.ChatID, err)
		}
		if i == 0 {
			d = &t
		}
		if t.Before(*d) {
			d = &t
		}
	}

	return d, nil
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
	var convoid string
	var participant string
	users, err := GetUsers()
	if err != nil {
		return nil, fmt.Errorf("GetConversations (GetUsers) error: %+v\n", err)
	}
	for rows.Next() {
		err := rows.Scan(&convoid, &participant)
		if err != nil {
			return conversations, fmt.Errorf("GetConversations rows.Scan error: %+v\n", err)
		}
		if i := convoInConvos(convoid, conversations); i >= 0 {
			convo := conversations[i]
			if convo.ConvoID == convoid {
				for _, u := range users {
					if u.ID == participant {
						convo.Participants = append(convo.Participants, u)
					}
				}
			}
			conversations[i] = convo
		} else {
			user := User{}
			for _, u := range users {
				if u.ID == participant {
					user = u
				}
			}
			conversations = append(conversations, &Conversation{
				ConvoID:      convoid,
				Participants: []User{user},
			})
		}
	}
	err = rows.Err()
	if err != nil {
		return conversations, err
	}
	// fmt.Println(conversations)
	return conversations, nil
}

func convoInConvos(convoID string, convos []*Conversation) int {
	for i, c := range convos {
		if convoID == c.ConvoID {
			return i
		}
	}
	return -1
}

func GetChats() ([]Chat, error) {
	chats := []Chat{}
	rows, err := DB.Query(`SELECT * FROM chats`)
	if err != nil {
		return chats, fmt.Errorf("GetChats DB Query error: %+v\n", err)
	}
	var convoid string
	var chatid string
	var sender string
	var date string
	var body string
	for rows.Next() {
		err := rows.Scan(&convoid, &chatid, &sender, &date, &body)
		if err != nil {
			return chats, fmt.Errorf("GetChats rows.Scan error: %+v\n", err)
		}
		chats = append(chats, Chat{
			ConvoID: convoid,
			ChatID:  chatid,
			Sender: User{
				ID: sender,
			},
			Date: date,
			Body: body,
		})
	}
	err = rows.Err()
	if err != nil {
		return chats, err
	}
	return populateUsersForChats(chats)
}

func populateChatsForConversation(conversations ...*Conversation) error {
	chats, err := GetChats()
	if err != nil {
		return err
	}
	for i, convo := range conversations {
		cts := []Chat{}
		for _, c := range chats {
			if c.ConvoID == convo.ConvoID {
				cts = append(cts, c)
			}
		}
		conversations[i].Chats = cts
	}
	return nil
}

func GetPopulatedConversations(conversations []*Conversation) ([]*Conversation, error) {
	convos, err := GetConversations()
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedConversation (GetConversations) error: %+v\n", err)
	}
	if conversations == nil {
		conversations = convos
	} else {
		for i, c := range conversations {
			for _, co := range convos {
				if c.ConvoID == co.ConvoID {
					conversations[i] = co
				}
			}
		}
	}
	if err := populateChatsForConversation(conversations...); err != nil {
		return nil, fmt.Errorf("GetPopulatedConversation (populateChatForConversation) error: %+v\n", err)
	}
	// b, err := json.Marshal(conversations)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(b))
	return conversations, nil
}

func populateUsersForChats(chats []Chat) ([]Chat, error) {
	users, err := GetUsers()
	if err != nil {
		return nil, fmt.Errorf("populateUsersForChats (GetUsers) error: %+v\n", err)
	}
	outChats := []Chat{}
	for _, chat := range chats {
		newChat := chat
		for _, user := range users {
			if chat.Sender.ID == user.ID {
				newChat.Sender = user
			}
		}
		outChats = append(outChats, newChat)
	}
	return outChats, nil
}

func FilterChatsForConvo(convoID string, chats []Chat) []Chat {
	out := []Chat{}
	for _, c := range chats {
		if convoID == c.ConvoID {
			out = append(out, c)
		}
	}
	// fmt.Printf("FilterChatsForConvo Chats: %+v\n", chats)
	return out
}

func ParticipantsToIds(users []User) []string {
	result := []string{}
	for _, item := range users {
		result = append(result, item.ID)
	}
	return result
}

func GetConvoID(participantIDs []string, conversations []*Conversation) (string, error) {
	v := map[string][]string{}
	for _, convo := range conversations {
		v[convo.ConvoID] = ParticipantsToIds(convo.Participants)
	}
	for convoID, vv := range v {
		if reflect.DeepEqual(participantIDs, vv) {
			return convoID, nil
		}
	}
	return "", nil
}

var fieldsThatNeverConflict map[string]bool = map[string]bool{"Password": true, "Gender": true}

// (u User) compare diffs two user structs and returns a slice of conflicting keys.
func (u User) Compare(t User) []string {
	o := []string{}
	uv, tv := reflect.ValueOf(u), reflect.ValueOf(t)

	for i := 0; i < uv.NumField(); i++ {
		if _, ok := fieldsThatNeverConflict[uv.Type().Field(i).Name]; ok {
			continue
		}
		if uv.Field(i).String() == tv.Field(i).String() {
			o = append(o, uv.Type().Field(i).Name)
		}
	}

	return o
}

func CreatePost(post *Post) (*Post, error) {
	stmt, err := DB.Prepare("INSERT INTO posts (postID, nickname, title, categories, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	// fmt.Println(post.Nickname)
	if err != nil {
		return nil, fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	if post.PostID == "" {
		post.PostID = uuid.NewV4().String()
	}
	// TODO: remove placeholder nickname once login/sessions are working, and hook up the real user who is logged in
	if post.Nickname == "" {
		fmt.Println("post.Nickname is empty")
		// post.Nickname = "Cassidy"
	}
	_, err = stmt.Exec(post.PostID, post.Nickname, post.Title, post.Categories, post.Body)
	if err != nil {
		return nil, fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return post, err
}

func CreateComment(comment Comment) (Comment, error) {
	stmt, err := DB.Prepare("INSERT INTO comments (commentID, postID, nickname, body) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return comment, fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	if comment.CommentID == "" {
		comment.CommentID = uuid.NewV4().String()
	}
	// TODO: remove placeholder nickname once login/sessions are working
	if comment.Nickname == "" {
		fmt.Println("comment.Nickname is empty")
		// comment.Nickname = "Cassidy"
	}
	_, err = stmt.Exec(comment.CommentID, comment.PostID, comment.Nickname, comment.Body)
	if err != nil {
		return comment, fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	return comment, err
}

func CreateUser(user User) (User, error) {
	stmt, err := DB.Prepare("INSERT INTO users (ID, nickname, age, gender, firstname, lastname, email, password) VALUES (?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err != nil {
		return user, fmt.Errorf("Create User DB Prepare error: %+v\n", err)
	}
	if user.ID == "" {
		user.ID = uuid.NewV4().String()
	}
	_, err = stmt.Exec(user.ID, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return user, fmt.Errorf("Create User Exec error: %+v\n", err)
	}
	return user, err
}

func GetUserByNickname(user User) (User, error) {
	if user.Nickname == "" {
		return User{}, fmt.Errorf("user must contain a nickname")
	}

	users, err := GetUsers()
	if err != nil {
		return User{}, fmt.Errorf("unable to get users for filtering user: %w", err)
	}

	for _, u := range users {
		if user.Nickname == u.Nickname {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("unable to find user for '%s'", user.Nickname)
}

const SessionCookieName string = "forum-session"

func CreateSession(user User) (*http.Cookie, *Session, error) {
	sessionID := uuid.NewV4().String()
	expiration := time.Now().Add(1 * time.Hour)
	cookie := &http.Cookie{Name: SessionCookieName, Value: sessionID, Expires: expiration, SameSite: http.SameSiteLaxMode}

	rows, err := DB.Prepare(`INSERT or REPLACE INTO sessions(sessionID, userID, expiryTime) VALUES (?,?,?);`)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to prepare insertion into sessions: %w", err)
	}
	defer rows.Close()
	_, err = rows.Exec(sessionID, user.ID, expiration)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to execute sessions insertion statement: %w", err)
	}
	return cookie, &Session{SessionID: sessionID, UserID: user.ID, ExpiryTime: expiration.String()}, nil
}
