package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbFile = os.Getenv("TODO_DBFILE")

// DB init or create if not exist
func initDB() {
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		fmt.Printf("Creating database...")
		db, err := sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {

			}
		}(db)

		query := `
			CREATE TABLE scheduler (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			date TEXT,
    			title TEXT NOT NULL,
    			comment TEXT,
    			repeat TEXT
   			);
   			CREATE INDEX idx_date ON scheduler(date);
		`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database created successfully!")
	}
}
