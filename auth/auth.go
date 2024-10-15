package auth

import (
	"net/http"

	"github.com/deeramster/go_final_project/config"
)

// Middleware для проверки авторизации
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := config.AppConfig.Password
		if pass == "" {
			next(w, r) // Если пароль не установлен, пропускаем без авторизации
			return
		}

		// Получаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value != pass {
			// Возвращаем ошибку аутентификации в формате JSON
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
