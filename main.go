package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)
var sqliteDatabase *sql.DB
func main() {
//Open the database SQLite file and create the database table
	database, err1 := sql.Open("sqlite3", "sqlite-database.db")
	sqliteDatabase = database
if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	defer sqliteDatabase.Close()

	myhttp := http.NewServeMux()
	fs := http.FileServer(http.Dir("./."))
	myhttp.Handle("/", http.StripPrefix("", fs))

myhttp.HandleFunc("/home", mainHandler)
	myhttp.HandleFunc("/chat", chatSocketCreate)
	myhttp.HandleFunc("/content", contentSocketCreate)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", myhttp)

	
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("index.html"))
	if err := tpl.Execute(w, Person); err != nil {
		http.Error(w, "No such file or directory: Internal Server Error 500", http.StatusInternalServerError)
	}
}