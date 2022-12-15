package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"real-time-forum/pkg/authentication"
	"real-time-forum/pkg/database"
	socket "real-time-forum/pkg/websockets"

	_ "github.com/mattn/go-sqlite3"
)

const databaseFilePath string = "sqlite-database.db"

func init() {
	// set dev mode:
	// `DEV=true go run .`
	// set dev mode for all `go run`'s in this terminal:
	// `export DEV=true` - now run go run
	devMode := false
	if _, ok := os.LookupEnv("DEV"); ok {
		devMode = true
	}
	database.InitialiseDB(databaseFilePath, devMode)
}

// func register(w http.ResponseWriter, r *http.Request){
// 	r.ParseForm()
// 	fmt.Println(r.Form)
// 	if err := tpl.ExecuteTemplate(sb, "login.template", nil); err != nil {
// 			return fmt.Errorf("loginExecuteTemplate error: %+v\n", err)
// 		}
// }
func main() {
	defer database.DB.Close()
	myhttp := http.NewServeMux()
	myhttp.HandleFunc("/", serveHomePage)
	// when adding a new websocket endpoint make sure to update the switch case in the websocket connection function to account for it
	myhttp.HandleFunc("/chat", socket.SocketCreate)
	myhttp.HandleFunc("/content", socket.SocketCreate)
	myhttp.HandleFunc("/post", socket.SocketCreate)
	myhttp.HandleFunc("/presence", socket.SocketCreate)
	myhttp.HandleFunc("/comment", socket.SocketCreate)
	myhttp.HandleFunc("/register", authentication.RegistrationHandler)
	myhttp.HandleFunc("/login", authentication.LoginHandler)
	fs := http.FileServer(http.Dir("./static"))
	myhttp.Handle("/static", fs)
	// myhttp.HandleFunc("/home", mainHandler)

	// daemonised functions
	// go socket.BroadcastPresences()

	fmt.Println("http://localhost:8080")
	err := http.ListenAndServe(":8080", myhttp)
	if err != nil {
		log.Fatal(err)
	}
}

// func mainHandler(w http.ResponseWriter, r *http.Request) {
// 	tpl := template.Must(template.ParseGlob("index.html"))
// 	if err := tpl.Execute(w, auth.Person); err != nil {
// 		http.Error(w, "No such file or directory: Internal Server Error 500", http.StatusInternalServerError)
// 	}
// }

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "index.html" {
		tpl, err := template.ParseGlob("templates/*")
		if err != nil {
			fmt.Printf("unable to parse templates: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !authentication.RequestIsLoggedIn(r) {
			if err := tpl.ExecuteTemplate(os.Stdout, "indexWithLogin.template", nil); err != nil {
				fmt.Printf("unable to render indexWithLogin.template: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := tpl.ExecuteTemplate(w, "indexWithLogin.template", nil); err != nil {
				fmt.Printf("unable to render indexWithLogin.template: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if err := tpl.ExecuteTemplate(w, "index.template", nil); err != nil {
			fmt.Printf("unable to render index.template: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	return
}
