package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/handlers"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
	ws "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub *ws.Hub

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	hub.Register(conn)
	defer func() {
		hub.Unregister(conn)
	}()

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		log.Printf("WebSocket message: %v", msg)

		switch msg["type"] {
		case "JOIN_BOARD":
			log.Printf("User joined board: %v", msg["boardId"])
			conn.WriteJSON(map[string]interface{}{
				"type": "JOINED_BOARD",
				"boardId": msg["boardId"],
			})

		case "PING":
			conn.WriteJSON(map[string]interface{}{
				"type": "PONG",
				"data": "pong",
			})

		default:
			response, _ := json.Marshal(msg)
			hub.BroadcastMessage(response)
		}
	}
}

func main() {
	hub = ws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	connStr := "host=db-1 port=5432 user=user password=password dbname=TaskflowDB sslmode=disable"
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ur := repository.NewUserRepoImpl(conn)
	br := repository.NewBoardRepoImpl(conn)
	lr := repository.NewListRepoImpl(conn)
	cr := repository.NewCardRepoImpl(conn)

	as := services.NewAuthService(ur)
	bs := services.NewBoardService(br, lr, cr, ur)
	ls := services.NewListService(lr, br, cr)
	cs := services.NewCardService(cr, lr, br)

	ah := handlers.NewAuthHandler(as)
	bh := handlers.NewBoardHandler(bs, as)
	lh := handlers.NewListHandler(ls, as, hub)
	ch := handlers.NewCardHandler(cs, as, hub)

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

	mux.HandleFunc("/ws", wsHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.208.203:8081", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(mux)

	log.Println("Сервер запущен на :8080 (HTTP + WebSocket)")
	http.ListenAndServe(":8080", handler)
}