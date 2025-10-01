package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/handlers"
)

func main() {

	mux := http.NewServeMux()
	/*
		db, err := sql.Open("sqlite3", "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/database/SQLite/store.db")
		if err != nil {
			panic(err)
		}
		defer db.Close()
	*/
	h := handlers.NewHandler()
	mux.HandleFunc("/api/auth/register", h.Register)
	mux.HandleFunc("/api/auth/login", h.Login)
	mux.HandleFunc("/api/auth/me", h.Me)
	mux.HandleFunc("/api/boards", h.GetBoards)
	mux.HandleFunc("/api/auth/logout", h.Logout)
	mux.HandleFunc("/api/boards/{id}", h.BoardDelete)
	mux.HandleFunc("/api/boards/{boardId}/restore", h.BoardRestore)

	// Настройка CORS с помощью библиотеки
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.208.203:8081", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(mux)

	println(" Сервер запущен на :8080")
	http.ListenAndServe(":8080", handler)
}
