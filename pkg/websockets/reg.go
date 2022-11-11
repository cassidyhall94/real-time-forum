package websockets
import (
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
	// tpl, err := template.ParseGlob("templates/*")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	var data database.User
	r.ParseForm()
	fmt.Println(r.Form)
	data.Nickname = r.FormValue("nickname")
	// age, err := strconv.Atoi( r.FormValue("age")); if err !=nil {
	// 	log.Fatal(err)
	// }
	// data.Age = age
	data.Age = r.FormValue("age")
	data.Gender = r.FormValue("gender")
	data.FirstName = r.FormValue("fname")
	data.LastName = r.FormValue("lname")
	data.Email = r.FormValue("email")
	// password := r.FormValue("password")
	password := passwordHash(r.FormValue("password"))
	var pwMatch = checkPwHash(r.FormValue("password"), password)
	fmt.Println("pw match", pwMatch)
	data.Password = password
	// log.Println(data)
	CreateUser(data)
	// if err := tpl.ExecuteTemplate(w, "login.template", nil); err != nil {
	// 	 fmt.Errorf("loginExecuteTemplate error: %+v\n", err)
	// }
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

