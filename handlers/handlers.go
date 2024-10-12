package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/deeramster/go_final_project/dateutil"
	"github.com/deeramster/go_final_project/db"
	"net/http"
	"strconv"
	"time"
)

func HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
			return
		}

		// Проверка обязательного поля Title
		if task.Title == "" {
			http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
			return
		}

		// Если дата не указана, присваиваем сегодняшнюю дату
		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		// Проверяем формат даты задачи
		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "Invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
			return
		}

		// Получаем сегодняшнюю дату
		today := time.Now()
		todayFormatted := today.Format("20060102")
		todayDate, _ := time.Parse("20060102", todayFormatted)

		// Логируем даты для отладки
		fmt.Printf("Comparing task date: %s with today's date: %s\n", taskDate.Format("20060102"), todayFormatted)

		// Сравниваем дату задачи с сегодняшней датой
		if taskDate.Before(todayDate) {
			// Если дата меньше сегодняшней
			if task.Repeat == "" {
				// Если правило повторения не указано, устанавливаем сегодняшнюю дату
				task.Date = todayFormatted
			} else {
				// Вычисляем следующую дату с помощью функции NextDate
				nextDate, err := dateutil.NextDate(today, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error": "Error calculating next date: %v"}`, err), http.StatusBadRequest)
					return
				}
				task.Date = nextDate // Устанавливаем вычисленную дату
			}
		}

		// Добавление задачи в базу данных
		id, err := db.AddTaskToDB(task)
		if err != nil {
			http.Error(w, `{"error": "Failed to add task"}`, http.StatusInternalServerError)
			return
		}

		// Возвращаем ID созданной задачи в формате JSON
		json.NewEncoder(w).Encode(map[string]int64{"id": id})
	case "PUT":
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := db.UpdateTaskInDB(task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	case "DELETE":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := db.DeleteTaskFromDB(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []db.Task
	var err error

	if search != "" {
		// Checking and format date from this layout DD.MM.YYYY
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			tasks, err = db.SearchTasksByDate(parsedDate.Format("20060102"))
		} else {
			// Searching by title or comment
			tasks, err = db.SearchTasksByText(search)
		}
	} else {
		// Return all tasks if search empty
		tasks, err = db.GetTasksFromDB()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string][]db.Task{"tasks": tasks})
}

func HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, err := db.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	nextDate, err := dateutil.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.MarkTaskDoneInDB(id, nextDate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid 'now' date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	if dateStr == "" {
		http.Error(w, "Parameter 'date' is required", http.StatusBadRequest)
		return
	}

	nextDate, err := dateutil.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating next date: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}
