package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

var DB *sql.DB

func InitialiseDB(path string, insertPlaceholders bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err.Error())
		}
		if insertPlaceholders {
			defer insertPlaceholdersInDB()
		}
		file.Close()
	}
	//Open the database SQLite file and create the database table
	sqliteDatabase, err1 := sql.Open("sqlite3", path)
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	//Defer the close
	// defer sqliteDatabase.Close()
	//Create the database for each user
	_, errTbl := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	TEXT,
			"email" 	TEXT UNIQUE,
			"username"	TEXT UNIQUE,
			"password"	TEXT,
			"nickname" TEXT
		);
	`)
	if errTbl != nil {
		fmt.Println("USER TABLE ERROR")
		log.Fatal(errTbl.Error())
	}

	//Create the posts table
	_, errPosts := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"postID"	TEXT,
			"userName"	TEXT NOT NULL,
			"category"	TEXT ,
			"title" TEXT,
			"postContent" TEXT
		);
	`)
	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}
	//Create the cookies table
	_, errCookie := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "cookies" (
			"name"	TEXT,
			"sessionID" 	TEXT UNIQUE
		);
	`)
	if errCookie != nil {
		fmt.Println("TABLE ERROR")
		log.Fatal(errTbl.Error())
	}

	//Create the table for each user
	_, errComments := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentID" TEXT,
			"postID"	TEXT,
			"username"	TEXT ,
			"commentText"	TEXT
		);
	`)
	if errComments != nil {
		fmt.Println("USER ERROR")
		log.Fatal(errTbl.Error())
	}
	//Create the database for each user
	_, errCategories := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "categories" (
			"postID"	TEXT,
			"Javascript"	INTEGER,
			"Go"	INTEGER,
			"Rust"	INTEGER		);
	`)
	if errCategories != nil {
		fmt.Println("CATEGORY TABLE ERROR")
		log.Fatal(errCategories.Error())
	}

	//Create the database for each user
	_, errChats := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "chats" (
			"chatID"	TEXT,
			"userID"	TEXT,
			"messageContent"	TEXT,
			"date"	INTEGER	);
	`)
	if errChats != nil {
		fmt.Println("CHAT TABLE ERROR")
		log.Fatal(errCategories.Error())
	}

	DB = sqliteDatabase
}

func insertPlaceholdersInDB() {
	queries := map[string]string{
		"add some dummy users": fmt.Sprintf(`INSERT INTO users values ("%s", "foo@bar.com", "foo", "s0fj489fhjsof", "bar")`, uuid.NewV4()),
	}

	for purpose, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			panic(fmt.Errorf("placeholder query '%s' failed with err: %w", purpose, err))
		} else {
			fmt.Println("successfully executed " + purpose)
		}
	}
}
