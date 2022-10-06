package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

	myhttp.HandleFunc("/chat", chatSocketCreate)
	myhttp.HandleFunc("/content", contentSocketCreate)
	myhttp.HandleFunc("/post", postSocketCreate)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", myhttp)
}
