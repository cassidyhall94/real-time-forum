package main

import (
	"database/sql"
	"fmt"
	"net/http"
	socket "real-time-forum/pkg/websockets"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase *sql.DB

func main() {
	database, err1 := sql.Open("sqlite3", "sqlite-database.db")
	sqliteDatabase = database
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	defer sqliteDatabase.Close()
	myhttp := http.NewServeMux()
	fs := http.FileServer(http.Dir("./."))
	myhttp.Handle("/", http.StripPrefix("", fs))

	// when adding a new websocket endpoint make sure to update the switch case in the websocket connection function to account for it
	myhttp.HandleFunc("/chat", socket.SocketCreate)
	myhttp.HandleFunc("/content", socket.SocketCreate)
	myhttp.HandleFunc("/post", socket.SocketCreate)
	// myhttp.HandleFunc("/home", mainHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", myhttp)

}

// func mainHandler(w http.ResponseWriter, r *http.Request) {
// 	tpl := template.Must(template.ParseGlob("index.html"))
// 	if err := tpl.Execute(w, auth.Person); err != nil {
// 		http.Error(w, "No such file or directory: Internal Server Error 500", http.StatusInternalServerError)
// 	}
// }
