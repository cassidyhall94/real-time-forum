package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	fmt.Println("SQL Databasefile created")
	sqliteDatabase, err1 := sql.Open("sqlite3", "sqlite-database.db")
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	defer sqliteDatabase.Close()
	_, errTbl := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	TEXT,
			"email" 	TEXT UNIQUE,
			"username"	TEXT UNIQUE,
			"password"	TEXT 
			"nickname"  TEXT
		);
	`)
	if errTbl != nil {
		fmt.Println("USER TABLE ERROR")
		log.Fatal(errTbl.Error())
	}
	_, errPosts := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"postID"	TEXT,
			"userName"	TEXT NOT NULL,
			"category"	TEXT ,
			"likes" INTEGER,
			"dislikes" INTEGER,
			"title" TEXT,
			"post" TEXT
		);
	`)
	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}
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

	_, errComments := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentID" TEXT,
			"postID"	TEXT,
			"username"	TEXT ,
			"commentText"	TEXT,
			
		);
	`)
	if errComments != nil {
		fmt.Println("USER ERROR")
		log.Fatal(errTbl.Error())
	}
	_, errCategories := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "categories" (
			"postID"	TEXT,
			"Javascript"	INTEGER,
			"Go"	INTEGER,
			"Rust"	INTEGER		);
	`)
	if errCategories != nil {
		fmt.Println("Creating Category table ERROR")
		log.Fatal(errCategories.Error())
	}

	_, errChats := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "chats" (
			"chatID"	TEXT,
			"userID"	INTEGER,
			"messageContent"	TEXT,
			"date"	INTEGER		);
	`)
	if errChats != nil {
		fmt.Println("Creating Chat table ERROR")
		log.Fatal(errCategories.Error())
	}
}
