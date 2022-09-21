package main

import (
	"database/sql"
	"fmt"
	database "forum/lib-database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	// "forum/database"

	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{} // use default options

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()
	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error during message writing:", err)
			break
		}
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/stylesheet.css")
}

func main() {

	// key:= "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http://127.0.0.1:3000"
	db, err := sql.Open("sqlite3", "./database/userdata.db")
	if err != nil {
		fmt.Printf("main (sql.Open) error: %+v\n", err)
		os.Exit(1)
	}
	db.Exec("PRAGMA journal_mode=WAL;")
	go func(db *sql.DB) {
		for {
			time.Sleep(time.Second * 10)
			// fmt.Println("Checkpointing WAL")
			_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Println("Checkpointed WAL")
		}
	}(db)
	defer db.Exec("VACCUM;")
	data, err := database.Connect(db)
	if err != nil {
		fmt.Printf("main func ConnectTables error: %+v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/socket", socketHandler)
	mux.HandleFunc("/", data.Handler)
	mux.HandleFunc("/category/", data.CategoryDump)
	mux.HandleFunc("/categoryg/", data.CategoryDump)
	mux.HandleFunc("/threadg/", data.ThreadGuest)
	mux.HandleFunc("/thread/", data.Threads)
	mux.HandleFunc("/category/stylesheet", cssHandler)
	mux.HandleFunc("/threadg/stylesheet", cssHandler)
	mux.HandleFunc("/thread/stylesheet", cssHandler)
	mux.HandleFunc("/categoryg/stylesheet", cssHandler)

	fmt.Println("Starting server at port 8080:\n http://localhost:8080/login")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(500, "500 Internal server error:", err)
		fmt.Printf("main ListenAndServe error: %+v\n", err)
	}

}
