package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() {
	// Получаем путь к базе данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Используем полный путь к файлу базы данных
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal("Cannot find executable path:", err)
	}
	dbFilePath := filepath.Join(filepath.Dir(executablePath), dbFile)
	fmt.Println(dbFilePath)

	// Is exist
	_, err = os.Stat(dbFilePath)
	if os.IsNotExist(err) {
		// If not exist, create new
		fmt.Println("Creating database...")

		db, err := sql.Open("sqlite3", dbFilePath)
		if err != nil {
			log.Fatal("Error opening DB:", err)
		}
		defer db.Close()

		// Create database and index
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
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Error creating table:", err)
		}

		fmt.Println("Database created successfully!")
	} else {
		fmt.Println("Database already exists.")
	}
}
