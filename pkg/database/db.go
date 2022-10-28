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
		fmt.Println(insertPlaceholders)
		if insertPlaceholders {
			defer insertPlaceholdersInDB()
		}
		file.Close()
	}
	// Open the database SQLite file and create the database table
	sqliteDatabase, err1 := sql.Open("sqlite3", path)
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	// Defer the close
	// defer sqliteDatabase.Close()
	// Create the database for each user
	_, errTbl := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	TEXT UNIQUE,
			"nickname" TEXT,
			"age" TEXT,
			"gender" TEXT,
			"firstname" TEXT,
			"lastname" TEXT,
			"email" 	TEXT UNIQUE,
			"password"	TEXT UNIQUE
		);
	`)
	if errTbl != nil {
		fmt.Println("USER TABLE ERROR")
		log.Fatal(errTbl.Error())
	}

	// Create the posts table
	_, errPosts := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"postID"	TEXT UNIQUE,
			"nickname"		TEXT,
			"categories"	TEXT,
			"title" TEXT,
			"body" TEXT
		);
	`)
	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}
	// Create the cookies table
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

	// Create the table for each user
	_, errComments := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentID" TEXT,
			"postID"	TEXT,
			"nickname"	TEXT,
			"body"	TEXT
		);
	`)
	if errComments != nil {
		fmt.Println("USER ERROR")
		log.Fatal(errTbl.Error())
	}
	// Create the database for each user
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

	// Create the database for each user
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
		"fake user 1": fmt.Sprintf(`INSERT INTO users values ("%s", "bar", "22", "female", "foodie", "barry", "foo@bar.com", "s0fj489fhjsof")`, uuid.NewV4()),

		"fake user 2": fmt.Sprintf(`INSERT INTO users values ("%s", "foo", "30", "male", "barry", "fool", "bar@foo.com", "03444f89fsof")`, uuid.NewV4()),

		"fake post 1": fmt.Sprintf(`INSERT INTO posts values ("9b4bc963-ecb2-4767-a79b-b09cd102ce4a", "bar", "golang", "Best Coding Language ever", "Golang is really the best!")`),

		"fake post 2": fmt.Sprintf(`INSERT INTO posts values ("16f94e48-82bc-4884-96b3-c847d37f069c", "foo", "javascript", "I love Javascript!", "JS is really neat!")`),

		"fake comment 1": fmt.Sprintf(`INSERT INTO comments values ("49f89e2f-4d7d-4b03-beb6-8def55652d4a", "9b4bc963-ecb2-4767-a79b-b09cd102ce4a", "Cassidy", "I like it too!")`),
		
		"fake comment 2": fmt.Sprintf(`INSERT INTO comments values ("fbbd419a-e40f-49d5-867a-afa328127cbb", "16f94e48-82bc-4884-96b3-c847d37f069c", "Jeff", "Thanks for this post!")`),
	}

	for purpose, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			panic(fmt.Errorf("placeholder query '%s' failed with err: %w", purpose, err))
		} else {
			fmt.Println("successfully executed " + purpose)
		}
	}
}
