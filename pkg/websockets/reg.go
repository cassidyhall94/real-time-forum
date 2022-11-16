package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"

	"golang.org/x/crypto/bcrypt"
)

type userData struct {
	id        string
	nickname  string
	age       int
	gender    string
	firstName string
	lastName  string
	email     string
	password  string
}

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

	CreateUser(data)
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

func Login(w http.ResponseWriter, r *http.Request) {

	// upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }
	// con, _ := upgrader.Upgrade(w, r, nil)


	var user database.Login

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	var users []database.User
	

	//selects nickname and password from user database
	rows, err := database.DB.Query(`SELECT nickname, password,loggedin FROM users`)
	if err != nil {
		log.Println(err)
	}

	var nickname string
	var password string
	var loggedin string


	for rows.Next() {
		err := rows.Scan(&nickname, &password, &loggedin)
		if err != nil {
			log.Fatal(err)
		}

		// compares data with front end, if user nick match, checks pw if match stores value
		if user.Nickname == nickname {
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
	fmt.Println(users)
		//checks len again to stop panic err && updates user logged in to true in DB
		if len(users)>0 && users[0].LoggedIn =="true" {
			var loggedin= "true"
			UpdateUser(user.Nickname, loggedin)
			
		}
		//sends data to js front end
	 json.NewEncoder(w).Encode(users)

}

func UpdateUser(nickname, loggedin string){
	fmt.Println(nickname)
	
	stmt,err:= database.DB.Prepare(`UPDATE "users" SET "loggedin" = ? WHERE "nickname" = ?`);if err !=nil {
		log.Println(err)
	}
	stmt.Exec(loggedin, nickname)
}

func Logout(w http.ResponseWriter, r * http.Request){
	var name = "test"
	var loggedin = "false"
	UpdateUser(name, loggedin)
}
