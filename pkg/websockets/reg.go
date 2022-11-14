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

	err := json.NewDecoder(r.Body).Decode(&data); if err !=nil {
		log.Println(err)
		
	}
	fmt.Println(data)

	
	data.Password =passwordHash(data.Password)
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
