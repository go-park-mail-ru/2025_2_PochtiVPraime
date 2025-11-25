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

	/*
		// Накатить миграции
		if err := goose.Up(conn.DB, "../../db/migrations"); err != nil {
			log.Fatal(err)
		}
	*/
	//repository
	ur := repository.NewUserRepoImpl(conn)
	br := repository.NewBoardRepoImpl(conn)
	lr := repository.NewListRepoImpl(conn)
	cr := repository.NewCardRepoImpl(conn)
	bmr := repository.NewBoardMemberRepoImpl(conn)
	cmr := repository.NewCardMemberRepository(conn)
	clr := repository.NewChecklistRepository(conn)
	clpr := repository.NewChecklistPointRepository(conn)
	comr := repository.NewCommentRepository(conn)
	//sr := repository.NewSupportRepoImpl(conn)

	//services
	as := services.NewAuthService(ur)
	bs := services.NewBoardService(br, lr, cr, ur)
	ls := services.NewListService(lr, br, cr)
	cs := services.NewCardService(cr, lr, br)
	bms := services.NewBoardMemberService(bmr)
	cms := services.NewCardMemberService(cmr)
	cls := services.NewChecklistService(clr, clpr)
	clps := services.NewChecklistPointService(clpr)
	coms := services.NewCommentService(comr)
	//ss := services.NewSupportService(sr)

	//handlers
	ah := handlers.NewAuthHandler(as)
	bh := handlers.NewBoardHandler(bs, as)
	lh := handlers.NewListHandler(ls, as)
	ch := handlers.NewCardHandler(cs, as)
	bmh := handlers.NewBoardMemberHandler(bms)
	cmh := handlers.NewCardMemberHandler(cms, as)
	clh := handlers.NewChecklistHandler(cls, as)
	clph := handlers.NewChecklistPointHandler(clps, as)
	comh := handlers.NewCommentHandler(coms, as)
	//sh := handlers.NewSupportHandler(ss, as)

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
	mux.HandleFunc("/api/boards/{boardId}/boardMembers", bmh.BoardMember)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/cardMembers", cmh.CardMember)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/checklists", clh.GetOrCreateChecklists)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/checklist/{checklistId}", clh.Checklist)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/checklist/{checklistId}/points", clph.GetOrCreateChecklistsPoints)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/checklist/{checklistId}/point/{pointId}", clph.ChecklistPoint)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/comments", comh.Comments)
	mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/comment/{commentId}", comh.Comments)

	//mux.HandleFunc("/api/board/{boardId}/list/{listId}/task/{taskId}/checklist",)

	//mux.HandleFunc("/api/forms", sh.CreateOrGetForms)
	//mux.HandleFunc("/api/forms/{formId}", sh.DeleteOrGetForm)
	//mux.HandleFunc("/api/forms/statistic", sh.GetAllSupportForms)

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
