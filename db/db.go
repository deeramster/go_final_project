package db

import (
	"database/sql"
	"log"

	"github.com/deeramster/go_final_project/config"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB opens a connection to the database
func InitDB() {
	var err error
	if db == nil {
		db, err = sql.Open("sqlite3", config.AppConfig.DBFile)
		if err != nil {
			log.Fatal("Error opening DB:", err)
		}
		// Initialize the database schema if necessary
		createDatabaseIfNotExist()
	}
}

// createDatabaseIfNotExist ensures that the database and tables are created
func createDatabaseIfNotExist() {
	query := `
  CREATE TABLE IF NOT EXISTS scheduler (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   date TEXT,
   title TEXT NOT NULL,
   comment TEXT,
   repeat TEXT
  );
  CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
 `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

// GetDB returns the active database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection when the application stops
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			return
		}
	}
}
