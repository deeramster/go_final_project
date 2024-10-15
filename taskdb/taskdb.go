package taskdb

import (
	"github.com/deeramster/go_final_project/db"
	"github.com/deeramster/go_final_project/models"
)

const maxTasksReturned = 50

// AddTaskToDB adds a new task to the database
func AddTaskToDB(task models.Task) (int64, error) {
	query := "INSERT INTO scheduler (title, date, comment, repeat) VALUES (?, ?, ?, ?)"
	result, err := db.GetDB().Exec(query, task.Title, task.Date, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetTasksFromDB retrieves tasks from the database
func GetTasksFromDB() ([]models.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := db.GetDB().Query(query, maxTasksReturned)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
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

// GetTaskByID retrieves a task by its ID
func GetTaskByID(taskID int) (models.Task, error) {
	query := "SELECT id, title, date, comment, repeat FROM scheduler WHERE id = ?"
	var task models.Task
	err := db.GetDB().QueryRow(query, taskID).Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

// UpdateTaskInDB updates an existing task in the database
func UpdateTaskInDB(task models.Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	_, err := db.GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

// DeleteTaskFromDB deletes a task by its ID
func DeleteTaskFromDB(taskID int) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := db.GetDB().Exec(query, taskID)
	return err
}

// MarkTaskAsDone deletes a task or updates it if it's recurring
func MarkTaskAsDone(taskID int, nextDate string) error {
	if nextDate == "" {
		// Delete non-recurring task
		query := "DELETE FROM scheduler WHERE id = ?"
		_, err := db.GetDB().Exec(query, taskID)
		return err
	} else {
		// Update date for recurring task
		query := "UPDATE scheduler SET date = ? WHERE id = ?"
		_, err := db.GetDB().Exec(query, nextDate, taskID)
		return err
	}
}

// SearchTasksByDate retrieves tasks by date
func SearchTasksByDate(date string) ([]models.Task, error) {
	query := `
  SELECT id, date, title, comment, repeat
  FROM scheduler
  WHERE date = ?
  ORDER BY date
  LIMIT ?
 `
	rows, err := db.GetDB().Query(query, date, maxTasksReturned)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
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

// SearchTasksByText searches for tasks by text in title or comment
func SearchTasksByText(search string) ([]models.Task, error) {
	searchPattern := "%" + search + "%"
	query := `
  SELECT id, date, title, comment, repeat
  FROM scheduler
  WHERE title LIKE ? OR comment LIKE ?
  ORDER BY date
  LIMIT ?
 `
	rows, err := db.GetDB().Query(query, searchPattern, searchPattern, maxTasksReturned)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
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
