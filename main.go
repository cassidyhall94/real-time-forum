package main

import (
	"database/sql"
	"fmt"
	"net/http"
	auth "real-time-forum/pkg/authentication"
	socket "real-time-forum/pkg/websockets"
	"text/template"

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

	myhttp.HandleFunc("/chat", socket.ChatSocketCreate)
	myhttp.HandleFunc("/content", socket.ContentSocketCreate)
	myhttp.HandleFunc("/post", socket.PostSocketCreate)
	myhttp.HandleFunc("/home", mainHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", myhttp)

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("index.html"))
	if err := tpl.Execute(w, auth.Person); err != nil {
		http.Error(w, "No such file or directory: Internal Server Error 500", http.StatusInternalServerError)
	}
}
