package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title,omitempty" binding:"required"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func InitDB() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Is exist
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		// If not exist, create new
		fmt.Println("Creating database...")
		db, err := sql.Open("sqlite3", dbFile)
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

func AddTaskToDB(task Task) (int64, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetTasksFromDB() ([]Task, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(id int) (Task, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	var task Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, fmt.Errorf("task not found")
		}
		return Task{}, err
	}

	return task, nil
}

func MarkTaskDoneInDB(taskID int, nextDate string) error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if nextDate == "" {
		// Обычная задача, просто удаляем
		query := "DELETE FROM scheduler WHERE id = ?"
		_, err = db.Exec(query, taskID)
	} else {
		// Повторяющаяся задача, обновляем дату
		query := "UPDATE scheduler SET date = ? WHERE id = ?"
		_, err = db.Exec(query, nextDate, taskID)
	}
	return err
}

func SearchTasksByDate(date string) ([]Task, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
  		SELECT id, date, title, comment, repeat 
  		FROM scheduler 
  		WHERE date = ?
  		ORDER BY date
  		LIMIT 50
 	`
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func SearchTasksByText(search string) ([]Task, error) {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Подготовка поисковой строки для использования в SQL LIKE
	searchPattern := "%" + strings.ToLower(search) + "%"

	query := `
  		SELECT id, date, title, comment, repeat 
  		FROM scheduler 
  		WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ?
  		ORDER BY date
  		LIMIT 50
 	`

	rows, err := db.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func UpdateTaskInDB(task Task) error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	_, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

func DeleteTaskFromDB(taskID int) error {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM scheduler WHERE id = ?"
	_, err = db.Exec(query, taskID)
	return err
}
