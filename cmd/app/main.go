package main

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/handlers"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

func main() {

	mux := http.NewServeMux()
	connStr := "host=127.0.0.1 port=54320 user=user password=password dbname=TaskflowDB sslmode=disable"
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//repository
	ur := repository.NewUserRepoImpl(conn)
	br := repository.NewBoardRepoImpl(conn)
	lr := repository.NewListRepoImpl(conn)
	cr := repository.NewCardRepoImpl(conn)

	//services
	as := services.NewAuthService(ur)
	bs := services.NewBoardService(br, lr, cr, ur)

	//handlers
	ah := handlers.NewAuthHandler(as)
	bh := handlers.NewBoardHandler(bs, as)

	mux.HandleFunc("/api/auth/register", ah.Register)
	mux.HandleFunc("/api/auth/login", ah.Login)
	mux.HandleFunc("/api/auth/me", ah.Me)
	mux.HandleFunc("/api/boards", bh.GetBoards)
	mux.HandleFunc("/api/auth/logout", ah.Logout)
	mux.HandleFunc("/api/boards/{id}", bh.BoardDelete)
	mux.HandleFunc("/api/boards/{boardId}/restore", bh.BoardRestore)
	//mux.HandleFunc("/api/boards/{boardId}", bh.)
	//mux.HandleFunc("/api/boards/{boardId}/restore", bh.BoardRestore)
	//mux.HandleFunc("/api/boards/{boardId}/restore", bh.BoardRestore)

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
