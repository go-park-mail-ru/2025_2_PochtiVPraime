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
	connStr := "host=localhost port=54320 user=user password=password dbname=TaskflowDB sslmode=disable"
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
	sr := repository.NewSupportRepoImpl(conn)

	//services
	as := services.NewAuthService(ur)
	bs := services.NewBoardService(br, lr, cr, ur)
	ls := services.NewListService(lr, br, cr)
	cs := services.NewCardService(cr, lr, br)
	ss := services.NewSupportService(sr)

	//handlers
	ah := handlers.NewAuthHandler(as)
	bh := handlers.NewBoardHandler(bs, as)
	lh := handlers.NewListHandler(ls, as)
	ch := handlers.NewCardHandler(cs, as)
	sh := handlers.NewSupportHandler(ss, as)

	mux.HandleFunc("/api/auth/register", ah.Register)
	mux.HandleFunc("/api/user/profile", ah.UserUpdate)
	mux.HandleFunc("/api/user/password", ah.PasswordUpdate)
	mux.HandleFunc("/api/auth/login", ah.Login)
	mux.HandleFunc("/api/auth/me", ah.Me)
	mux.HandleFunc("/api/auth/logout", ah.Logout)
	mux.HandleFunc("/api/boards", bh.CreateOrGetBoards)
	mux.HandleFunc("/api/boards/{boardId}", bh.Board)
	mux.HandleFunc("/api/boards/{boardId}/restore", bh.BoardRestore)
	mux.HandleFunc("/api/boards/{boardId}/close", bh.ArchivedBoard)
	mux.HandleFunc("/api/board/{boardId}/lists", lh.CreateOrGetLists)
	mux.HandleFunc("/api/board/{boardId}/lists/{listId}", lh.List)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/tasks", ch.CreateOrGetCards)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}", ch.Card)
	mux.HandleFunc("/api/forms", sh.CreateOrGetForms)
	mux.HandleFunc("/api/forms/{formId}", sh.DeleteOrGetForm)
	mux.HandleFunc("/api/forms/statistic", sh.GetAllSupportForms)

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
