package main

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/handlers"
)

func main() {

	h := handlers.NewHandler()

	http.HandleFunc("/register", h.Register)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/get-boards", h.GetBoards)

	// TODO: Разобраться как и для каких ручек настроить CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	})

	println(" Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
