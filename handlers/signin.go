package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deeramster/go_final_project/config"
)

type Credentials struct {
	Password string `json:"password"`
}

func HandleSignIn(w http.ResponseWriter, r *http.Request) {
	// Если метод GET, возвращаем страницу логина
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./web/login.html")
		return
	}

	// Обработка POST-запроса для авторизации
	if r.Method == http.MethodPost {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		storedPassword := config.AppConfig.Password
		if storedPassword == "" {
			http.Error(w, `{"error": "Server is not configured"}`, http.StatusInternalServerError)
			return
		}

		if creds.Password != storedPassword {
			http.Error(w, `{"error": "Неверный пароль"}`, http.StatusUnauthorized)
			return
		}

		// Возвращаем токен в случае успешной аутентификации
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": creds.Password})
	}
}
