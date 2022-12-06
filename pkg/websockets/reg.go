package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// ***************************REGISTER**********************************************************8
// check if pasword meets criteria number length etc, if nickname is not taken
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting data")

	var data database.User

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)

	}

	fmt.Println(data)
	fmt.Println(data.LoggedIn)

	data.Password = passwordHash(data.Password)
	//  data.Password = checkPwHash(r.FormValue("password"), data.Password)

	var existingUsers, _ = database.GetUsers()
	var conflict = false

	//loop througn existing users to check if username or email is taken
	for i := 0; i < len(existingUsers); i++ {

		//if taken  breaks for loop and returns what value is taken and sets conflict to true
		if existingUsers[i].Nickname == data.Nickname {
			json.NewEncoder(w).Encode("username is taken")
			conflict = true
			break
		}
		if existingUsers[i].Email == data.Email {
			json.NewEncoder(w).Encode("email is taken")
			conflict = true
			break
		}

	}

	//if no conflicts adds user to db and resets conflict bool
	if !conflict {
		CreateUser(data)
		json.NewEncoder(w).Encode("registered")
		conflict = false
	}

}
func passwordHash(str string) string {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(str), 8)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPw)
}

func checkPwHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// *****************************LOGIN ***************************************
func Login(w http.ResponseWriter, r *http.Request) {

	var user database.Login

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	var users []database.User

	//selects nickname and password from user database
	rows, err := database.DB.Query(`SELECT nickname, password,loggedin, email FROM users`)
	if err != nil {
		log.Println(err)
	}

	var nickname string
	var password string
	var loggedin string
	var email string

	var store = sessions.NewCookieStore([]byte("secret-keys"))
	store.Options.SameSite = http.SameSiteLaxMode

	for rows.Next() {
		err := rows.Scan(&nickname, &password, &loggedin, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Nickname, nickname, email)

		// compares data with front end, if user nick match, checks pw if match stores value
		if user.Nickname == nickname || user.Nickname == email {
			if checkPwHash(user.Password, password) {
				users = append(users, database.User{
					Nickname: nickname,
					Password: password,
					LoggedIn: "true",
				})

			}

		}
	}

	//if len ==0, no matching user was found
	if len(users) == 0 {
		fmt.Println("pw mismatch")
	}
	// fmt.Println(users)

	//checks len again to stop panic err && updates user logged in to true in DB and creates cookie
	if len(users) > 0 && users[0].LoggedIn == "true" {
		// session, _ := store.Get(r, "session")
		// session.Values[nickname] = nickname
		// session.Save(r, w)
		// session.Options.SameSite = http.SameSiteLaxMode

		var loggedin = "true"
		UpdateUser(user.Nickname, loggedin)

		var cookieValue = uuid.NewV4()
		Cookie(w, r, user.Nickname, (cookieValue.String()))

	}
	//sends data to js front end
	json.NewEncoder(w).Encode(users)

}

//updates user table
func UpdateUser(nickname, loggedin string) {

	stmt, err := database.DB.Prepare(`UPDATE "users" SET "loggedin" = ? WHERE "nickname" = ?`)
	if err != nil {
		log.Println(err)
	}
	stmt.Exec(loggedin, nickname)
}

//logout user NOT WORKING YET
func Logout(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(r.Body)
	var usr struct{
		Nickname string `json:"nickname,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("current user logging out",usr)
	fmt.Println(usr.Nickname)
	usr.Nickname = strings.TrimSpace(usr.Nickname)
		fmt.Println(usr.Nickname)


	_, err = database.DB.Exec("DELETE FROM cookies where userName = '"+ usr.Nickname+ "'")
	if err !=nil {
		log.Fatal()
	}

	// cookie, err := r.Cookie("jerry")
	// if err !=nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println("logout cookie", cookie)

	// var name = "test"
	// var loggedin = "false"
	// UpdateUser(name, loggedin)
}

//creates cookie
func Cookie(w http.ResponseWriter, r *http.Request, Username string, id string) {

	expiration := time.Now().Add(1 * time.Hour)
	cookie := http.Cookie{Name: Username, Value: id, Expires: expiration, SameSite: http.SameSiteLaxMode}

	http.SetCookie(w, &cookie)

	rows, err := database.DB.Prepare(`INSERT INTO cookies(sessionID, userName, expiryTime) VALUES (?,?,?);`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	_, err = rows.Exec(id, Username, expiration)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(cookie)

}

func CheckCookies(w http.ResponseWriter, r *http.Request) {

}

// func makeSession() {
// }
