package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p PostFeed) MarshallJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (c Comment) MarshallJSON() ([]byte, error) {
	return json.Marshal(c)
}

// @TODO: error handling.
// login page.
func (data *Forum) LoginWeb(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	loggedIn := data.CheckCookie(w, r)
	// 🐈
	if loggedIn {
		http.Redirect(w, r, "/home", http.StatusFound)
	}

	// switch r.Method {
	// case "POST":
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("LoginWeb ParseForm error: %+v\n", err)
		return
	}

	var user User
	sessionToken := uuid.NewV4()
	expiresAt := time.Now().Add(720 * time.Second)

	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")

	var passwordHash string

	row := data.DB.QueryRow("SELECT password FROM people WHERE Username = ?", user.Username)
	err = row.Scan(&passwordHash)

	tpl := template.Must(template.ParseGlob("templates/*"))
	if err != nil {
		err := tpl.ExecuteTemplate(w, "login.html", "check username and password")
		if err != nil {
			fmt.Printf("LoginWeb ExecuteTemplate error: %+v\n", err)
			return
		}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))
	// returns nil on succcess
	if err == nil {
		err = data.CreateSession(Session{
			SessionID: sessionToken.String(),
			Username:  user.Username,
			Expiry:    expiresAt,
		})
		if err != nil {
			fmt.Printf("LoginWeb CreateSession error: %+v\n", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken.String(),
			Expires: expiresAt,
			// MaxAge:  2 * int(time.Hour),
		})
		// fmt.Println(data.GetSessions())
		// w.WriteHeader(200)
		http.Redirect(w, r, "/home", http.StatusFound)
		// data.HomePage(w, r)
	} else {
		fmt.Println("invalid credentials")
		err := tpl.ExecuteTemplate(w, "login.html", "check username and password")
		if err != nil {
			fmt.Printf("LoginWeb ExecuteTemplate error: %+v\n", err)
			return
		}
	}
}

func (data *Forum) GetSignupPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	err := tpl.ExecuteTemplate(w, "signup.html", nil)
	if err != nil {
		fmt.Printf("GetSignupPage Execute (signup.html) error: %+v\n", err)
	}
}

/*  1. check e-mail criteria
    2. check u.username criteria
	 3. check password criteria
	 4. check if u.username is already exists in database
	 5. create bcrypt hash from password
	 6. insert u.username and password hash in database
*/
func (data *Forum) SignUpUser(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*"))
	err := r.ParseForm() // parses sign up form to fetch needed information
	if err != nil {
		fmt.Printf("SignUpUser ParseForm error: %+v\n", err)
	}

	var user User

	user.Email = r.FormValue("email")
	// check if e-mail is valid format
	isValidEmail := true

	if isValidEmail != strings.Contains(user.Email, "@") || isValidEmail != strings.Contains(user.Email, ".") { // checks if e-mail is valid by checking if it contains "@"
		isValidEmail = false
	}

	if !isValidEmail {
		fmt.Println("Email invalid")
	}

	user.Username = r.FormValue("username")
	// check if only alphanumerical numbers
	isAlphaNumeric := true

	for _, char := range user.Username {
		if unicode.IsLetter(char) && unicode.IsNumber(char) { // checks if character not a special character
			isAlphaNumeric = false
		}
	}
	// checks if name length meets criteria
	nameLength := (5 <= len(user.Username) && len(user.Username) <= 50)

	// fmt.Println(nameLength)

	// check pw criteria
	user.Password = r.FormValue("password")

	// fmt.Println(user)
	var pwLower, pwUpper, pwNumber, pwSpace, pwLength bool
	pwSpace = false

	for _, char := range user.Password {
		switch {
		case unicode.IsLower(char):
			pwLower = true
		case unicode.IsUpper(char):
			pwUpper = true
		case unicode.IsNumber(char):
			pwNumber = true
		// case unicode.IsPunct(char) || unicode.IsSymbol(char):
		// 	pwSpecial = true
		case unicode.IsSpace(int32(char)):
			pwSpace = true
		}
	}
	minPwLength := 8
	maxPwLength := 30

	if minPwLength <= len(user.Password) && len(user.Password) <= maxPwLength {
		pwLength = true
	}

	if !pwLower || !pwUpper || !pwNumber || !pwLength || pwSpace || !isAlphaNumeric || !nameLength || !isValidEmail {
		err := tpl.ExecuteTemplate(w, "signup.html", "please check username, password and e-mail are valid")
		if err != nil {
			fmt.Printf("SignUpUser ExecuteTemplate signup.html error: %+v\n", err)
			return
		}
		return
	}

	row := data.DB.QueryRow("SELECT uuid FROM people where username =?", user.Username)
	var username string
	err2 := row.Scan(&username)
	if err2 != sql.ErrNoRows {
		fmt.Printf("sql scan row user error: %+v\n", err2)
		err3 := tpl.ExecuteTemplate(w, "signup.html", "username taken")
		if err3 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (username) error1: %+v\n", err3)
			return
		}
	}
	row = data.DB.QueryRow("SELECT uuid FROM people where email =?", user.Email)
	var userEmail string
	err = row.Scan(&userEmail)
	if err != sql.ErrNoRows {
		fmt.Printf("sql scan row email error: %+v\n", err)
		err2 := tpl.ExecuteTemplate(w, "signup.html", "e-mail in use")
		if err2 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (email) error2: %+v\n", err2)
			return
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err2 := tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		if err2 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (password) error:  %+v\n", err2)
		}
		return
	}

	sessionID := uuid.NewV4()
	err = data.CreateUser(User{
		Uuid:     sessionID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: string(passwordHash),
	})

	if err != nil {
		err4 := tpl.ExecuteTemplate(w, "signup.html", "there was an error registering account")
		if err4 != nil {
			fmt.Printf("SignUpUser ExecuteTemplate (password) error:  %+v\n", err4)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (data *Forum) GetUsernameFromSessionID(writer http.ResponseWriter, request *http.Request) string {
	c, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Printf("GetUsernameFromSessionID (Cookie) error: %+v\n", err)
			return ""
		}
	}

	sessionToken := c.Value
	a, err := data.GetSessions()
	if err != nil {
		fmt.Printf("GetUsernameFromSessionID (GetSessions) error: %+v\n", err)
		return ""
	}

	for _, sess := range a {
		// fmt.Println(sessionToken, " : ", sess.SessionID)
		if sessionToken == sess.SessionID {
			return sess.Username
		}
	}
	return ""
}

// check cookie
func (data *Forum) CheckCookie(writer http.ResponseWriter, request *http.Request) bool {
	c, err := request.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Printf("CheckCookie (Cookie) error: %+v\n", err)
			return false
		}
	}

	sessionToken := c.Value
	sessions, err := data.GetSessions()
	if err != nil {
		fmt.Printf("CheckCookie (GetSessions) error: %+v\n", err)
	}

	for _, sess := range sessions {
		if sessionToken == sess.SessionID {
			err := data.AssertUniqueSessionForUser(sess)
			if err != nil {
				fmt.Printf("could not determine unique session for user %+v; %+v", sess, err)
				return false
			}
			return true
		}
	}

	return false
}

func (data *Forum) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		fmt.Printf("Logout Cookie error: %+v\n", err)
		return
	} else if c.Value == "" {
		fmt.Printf("Cookie not found: %+v\n", err)
		return
	}

	sessionToken := c.Value
	var currentSession Session
	a, err := data.GetSessions()
	if err != nil {
		fmt.Printf("Logout GetSessions error: %+v\n", err)
	}

	for _, sess := range a {
		if sessionToken == sess.SessionID {
			currentSession = sess
			_, err := data.DB.Exec("DELETE FROM session where sessionID ='" + currentSession.SessionID + "'")
			if err != nil {
				fmt.Printf("Logout Exec error: %+v\n", err)
			}
		}
	}

	c = &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// home page.
func (data *Forum) HomePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	tpl := template.Must(template.ParseGlob("templates/*"))

	if err := request.ParseForm(); err != nil { // checks for errors parsing form
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("HomePage (ParseForm) error:  %+v\n", err)
		return
	}
	loggedIn := data.CheckCookie(writer, request)
	// 🐈
	if !loggedIn {
		data, err := data.GetPosts()
		if err != nil {
			fmt.Printf("HomePage (GetPosts) error: %+v\n", err)
		}
		err = tpl.ExecuteTemplate(writer, "guest.html", data)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (guest.html) error: %+v\n", err)
		}
		return
	} else {
		post, err := data.GetPosts()
		if err != nil {
			fmt.Printf("HomePage GetPosts (post) error: %+v\n", err)
		}
		var lastPost PostFeed
		if len(post) > 0 {
			lastPost = post[len(post)-1]
		}

		postCategory := request.FormValue("category")
		postCategory2 := request.FormValue("category2")
		// fmt.Println( postCategory2)

		if postCategory2 != "" {
			postCategory += " "
			postCategory += postCategory2
		}
		// fmt.Println(postCategory)
		// fmt.Println(postCategory)
		postTitle := request.FormValue("title")
		postContent := request.FormValue("content")

		postLikes := 0
		postDislikes := 0
		time := time.Now()
		postCreated := time.Format("01-02-2006 15:04")

		// checks session and selects the last one (the latest one)
		sess, err := data.GetSessions()
		if err != nil {
			fmt.Printf("HomePage GetSessions error: %+v\n", err)
		}
		currentSession := sess[len(sess)-1]
		user := currentSession.Username // fetches username from session

		type postSessionStruct struct {
			Post        []PostFeed
			UserSession Session
		}
		var postAndSession postSessionStruct

		postAndSession.UserSession = currentSession
		// checks if last post == current submit values to prevent duplicate posts
		if lastPost.Content == postContent {
			fmt.Println("duplicate")
			postAndSession.Post, err = data.GetPosts()
			if err != nil {
				fmt.Printf("HomePage GetPosts (lastPost.Content) error: %+v\n", err)
			}
			err := tpl.ExecuteTemplate(writer, "home.html", postAndSession)
			if err != nil {
				fmt.Printf("HomePage Execute (home.html) error: %+v\n", err)
			}
			return
		} else {
			// postAndSession.UserSession = data.GetSessions()[0]
			if postTitle != "" || postContent != "" {
				err := data.CreatePost(PostFeed{
					Username:  user,
					Title:     postTitle,
					Content:   postContent,
					Likes:     postLikes,
					Dislikes:  postDislikes,
					Category:  postCategory,
					CreatedAt: postCreated,
				})
				if err != nil {
					fmt.Printf("HomePage (CreatePost) items error: %+v\n", err)
					return
				}

				postAndSession.Post, err = data.GetPosts()
				if err != nil {
					fmt.Printf("HomePage (GetPosts) items error: %+v\n", err)
					return
				}

				err = tpl.ExecuteTemplate(writer, "home.html", postAndSession)
				if err != nil {
					fmt.Printf("HomePage ExecuteTemplate user homepage error: %+v\n", err)
					return
				}
				return

			}
		}
		data, err := data.GetPosts()
		postAndSession.Post = data
		if err != nil {
			fmt.Printf("HomePage (GetPosts) data error: %+v\n", err)
			return
		}
		err = tpl.ExecuteTemplate(writer, "home.html", postAndSession)
		if err != nil {
			fmt.Printf("HomePage ExecuteTemplate (home.html) error: %+v\n", err)
			return
		}
		return
	}
}

func (data *Forum) Handler(w http.ResponseWriter, r *http.Request) {
	// data.CheckCookie(w, r)

	switch r.URL.Path {
	// page handlers
	case "/stylesheet": // handle css
		http.ServeFile(w, r, "./templates/stylesheet.css")
	case "/":
		data.LoginWeb(w, r)
	case "/login":
		data.LoginWeb(w, r)
	case "/logout":
		data.Logout(w, r)
	case "/home":
		data.HomePage(w, r)
	case "/signup":
		data.GetSignupPage(w, r)
	case "/sign-up-form":
		data.SignUpUser(w, r)

		// handles images
	case "/cat":
		http.ServeFile(w, r, "./images/cat.jpg")
	case "/chicken":
		http.ServeFile(w, r, "./images/chicken.jpeg")
	case "/cow":
		http.ServeFile(w, r, "./images/cow.jpg")
	case "/hamster":
		http.ServeFile(w, r, "./images/hamster.jpg")
	case "/owl":
		http.ServeFile(w, r, "./images/owl.jpg")
	case "/panda":
		http.ServeFile(w, r, "./images/panda.jpg")
	case "/shark":
		http.ServeFile(w, r, "./images/shark.jpg")
	case "/doge":
		http.ServeFile(w, r, "./images/doge.jpg")
	case "/question":
		http.ServeFile(w, r, "./images/question.jpg")
	case "/finance":
		http.ServeFile(w, r, "./images/finance.jpg")
	case "/fitness":
		http.ServeFile(w, r, "./images/fitness.jpg")
	case "/health":
		http.ServeFile(w, r, "./images/health.jpg")
	case "/tech":
		http.ServeFile(w, r, "./images/tech.jpg")
	case "/travel":
		http.ServeFile(w, r, "./images/travel.jpg")
	}
}
