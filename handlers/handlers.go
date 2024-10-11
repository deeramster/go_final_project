package handlers

import (
	"encoding/json"
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
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		if task.Date == "" {
			task.Date = time.Now().Format("20060102")
		}

		id, err := db.AddTaskToDB(task)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]int64{"id": id})
	}
}

func HandleTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTasksFromDB()
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
