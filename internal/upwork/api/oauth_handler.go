package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func OAuthHandler() {
	r := chi.NewRouter()

	r.Get("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Код авторизации не был получен", http.StatusBadRequest)
			return
		}

		// Здесь код для обмена полученного кода на токен доступа
		// Вызов функции, которая обращается к серверу Upwork для получения токена
		fmt.Fprintf(w, "Код авторизации успешно получен: %s\n", code)
		// Добавить логику для обмена кода на токен и использование токена для доступа к API
	})

	// Запуск веб-сервера на порту 8080
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
