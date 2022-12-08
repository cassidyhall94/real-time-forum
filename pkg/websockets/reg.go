package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// *******************************************************CHECK IF SESSION********************************888
func CheckLoggedIn(r *http.Request) (bool, string){
	var cookies = getCookies(r)
	// get cookie
	// compare with db if value exists
	// if exists start websocket connection?
	
	if len(cookies) > 0 {
		fmt.Println(database.GetSessionsFromDB())
		// // sess = &[]database.Session{}
		var sess, _ = database.GetSessionsFromDB()
		for i := 0; i < len(sess); i++ {
			if cookies[0].Name == sess[i].UserName{
				return true, sess[i].UserName
			}
		}
	} 
	//returns false if cookie expired
return false, ""

}

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

	// var cookies = r.Cookies()
	var cookies = getCookies(r)

	if len(cookies) >= 1 {
		for i := 0; i < len(cookies); i++ {
			if cookies[i].Name != user.Nickname {

				DeleteCookie(w, cookies[i].Name)
				UpdateUser(cookies[i].Name, "false")
			}
		}
	}

	var users []database.User

	//selects nickname and password from user database
	rows, err := database.DB.Query(`SELECT ID, nickname, password,loggedin, email FROM users`)
	if err != nil {
		log.Println(err)
	}

	var nickname string
	var password string
	var loggedin string
	var email string
	var id string
	var matchID string

	// var store = sessions.NewCookieStore([]byte("secret-keys"))
	// store.Options.SameSite = http.SameSiteLaxMode

	for rows.Next() {
		err := rows.Scan(&id, &nickname, &password, &loggedin, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Nickname, nickname, email)

		// compares data with front end, if user nick match, checks pw if match stores value
		if user.Nickname == nickname || user.Nickname == email {
			if checkPwHash(user.Password, password) {
				matchID = id
				users = append(users, database.User{
					Nickname: nickname,
					Password: password,
					LoggedIn: "true",
				})
				user.Nickname = nickname

			}

		}
	}

	//if len ==0, no matching user was found
	if len(users) == 0 {
		fmt.Println("pw mismatch")
	}
	if len(users) > 0 && users[0].LoggedIn == "true" {

		var loggedin = "true"
		UpdateUser(user.Nickname, loggedin)

		// Cookie(w, r, user.Nickname, (cookieValue.String()))
		Cookie(w, r, user.Nickname, matchID)

	}
	//sends data to js front end
	json.NewEncoder(w).Encode(users)

}

//updates logged in status in user table
func UpdateUser(nickname, loggedin string) {

	stmt, err := database.DB.Prepare(`UPDATE "users" SET "loggedin" = ? WHERE "nickname" = ?`)
	if err != nil {
		log.Println(err)
	}
	stmt.Exec(loggedin, nickname)
}

// ********************************LOGOUT*********************************

func Logout(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(r.Body)
	var usr struct {
		Nickname string `json:"nickname,omitempty"`
	}

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		log.Println(err)
	}

	usr.Nickname = strings.TrimSpace(usr.Nickname)

	_, err = database.DB.Exec("DELETE FROM cookies where userName = '" + usr.Nickname + "'")
	if err != nil {
		log.Fatal()
	}

	DeleteCookie(w, usr.Nickname)
	UpdateUser(usr.Nickname, "false")

}

// *************************** COOKIES***********************
//creates cookie
func Cookie(w http.ResponseWriter,  r *http.Request, Username string, id string) {

	expiration := time.Now().Add(1 * time.Hour)
	cookie := http.Cookie{Name: Username, Value: id, Expires: expiration, SameSite: http.SameSiteLaxMode}

	http.SetCookie(w, &cookie)

	rows, err := database.DB.Prepare(`INSERT or REPLACE INTO cookies(sessionID, userName, expiryTime) VALUES (?,?,?);`)
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

func getCookies(r *http.Request) []*http.Cookie {
	return r.Cookies()
}

//delete cookie
func DeleteCookie(w http.ResponseWriter, username string) {

	cookie := &http.Cookie{
		Name:    username,
		Value:   "",
		Path:    "",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

}

func CheckExpiredCookies(){

}
