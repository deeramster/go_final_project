package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/deeramster/go_final_project/dateutil"
	"github.com/deeramster/go_final_project/models"
	"github.com/deeramster/go_final_project/taskdb"
)

func HandleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleTaskPost(w, r)
	case "PUT":
		handleTaskPut(w, r)
	case "GET":
		handleTaskGet(w, r)
	case "DELETE":
		handleTaskDelete(w, r)
	default:
		// Если метод не поддерживается, возвращаем 405 Method Not Allowed
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleTaskPost(w http.ResponseWriter, r *http.Request) {
	var task models.Task
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
	id, err := taskdb.AddTaskToDB(task)
	if err != nil {
		http.Error(w, `{"error": "Failed to add task"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем ID созданной задачи в формате JSON
	err = json.NewEncoder(w).Encode(map[string]int64{"id": id})
	if err != nil {
		return
	}
}

func handleTaskPut(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля Title
	if task.Title == "" {
		http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля ID
	if task.ID == "" { // Проверяем, что ID не пустая строка
		http.Error(w, `{"error": "Task ID is required"}`, http.StatusBadRequest)
		return
	}

	// Проверяем формат даты задачи
	if task.Date == "" {
		task.Date = time.Now().Format("20060102") // Присваиваем текущую дату, если не указана
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		http.Error(w, `{"error": "Invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
		return
	}

	// Дополнительные проверки для даты, если необходимо
	today := time.Now()
	todayFormatted := today.Format("20060102")
	todayDate, _ := time.Parse("20060102", todayFormatted)

	// Сравниваем дату задачи с сегодняшней датой
	if taskDate.Before(todayDate) {
		// Обработка ситуации, когда задача имеет дату в прошлом
		// Можно добавить дополнительные условия, если нужно
		http.Error(w, `{"error": "Task date cannot be in the past"}`, http.StatusBadRequest)
		return
	}

	// Обновление задачи в базе данных
	if err := taskdb.UpdateTaskInDB(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		return
	} // Возвращаем обновлённую задачу
}

func handleTaskGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// Преобразование idStr в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат идентификатора"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := taskdb.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
		return
	}

	// Формируем ответ в формате JSON
	response := map[string]string{
		"id":      task.ID,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	}

	// Устанавливаем заголовок Content-Type и возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
func handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Если ID отсутствует, возвращаем ошибку
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Преобразование idStr в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Если преобразование не удалось, возвращаем ошибку
		http.Error(w, `{"error": "Invalid ID format"}`, http.StatusBadRequest)
		return
	}

	// Удаление задачи из базы данных
	if err := taskdb.DeleteTaskFromDB(id); err != nil {
		// Если задача не найдена, возвращаем ошибку 404
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}
		// Если не удалось удалить задачу по другой причине, возвращаем ошибку
		http.Error(w, `{"error": "Failed to delete task"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем пустой JSON при успешном удалении
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Возвращаем статус 200 OK
	_, err = w.Write([]byte("{}"))
	if err != nil {
		return
	} // Пустой JSON
	return
}

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []models.Task
	var err error

	if search != "" {
		// Проверка и форматирование даты
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			tasks, err = taskdb.SearchTasksByDate(parsedDate.Format("20060102"))
		} else {
			// Поиск по заголовку или комментарию
			tasks, err = taskdb.SearchTasksByText(search)
		}
	} else {
		// Возвращаем все задачи, если строка поиска пустая
		tasks, err = taskdb.GetTasksFromDB()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	err = json.NewEncoder(w).Encode(map[string][]models.Task{"tasks": tasks})
	if err != nil {
		return
	}
}

func HandleTaskDone(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Преобразование idStr в int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID format"}`, http.StatusBadRequest)
		return
	}

	// Получаем задачу из базы данных
	task, err := taskdb.GetTaskByID(id)
	if err != nil {
		http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
		return
	}

	// Рассчитываем следующую дату, если задача повторяющаяся
	var nextDate string
	if task.Repeat != "" {
		nextDate, err = dateutil.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error": "Failed to calculate next date"}`, http.StatusBadRequest)
			return
		}
	}

	// Отмечаем задачу как выполненную
	if err := taskdb.MarkTaskAsDone(id, nextDate); err != nil {
		http.Error(w, `{"error": "Failed to mark task as done"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем статус 204 No Content (пустое тело)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Возвращаем статус 200 OK
	_, err = w.Write([]byte("{}"))
	if err != nil {
		return
	} // Пустой JSON
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
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}
}
