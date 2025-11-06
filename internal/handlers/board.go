package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

// Handler — обработчик HTTP-запросов
// Внедряет зависимости: BoardService и AuthService
type BoardHandler struct {
	BoardService services.BoardService
	AuthService  services.AuthService
}

// NewHandler — конструктор для Dependency Injection
// Создаёт и инициализирует все зависимости
func NewBoardHandler(boardService *services.BoardService, authService *services.AuthService) *BoardHandler {
	return &BoardHandler{
		BoardService: *boardService,
		AuthService:  *authService,
	}
}

// GetBoards — обрабатывает GET /api/boards
// --TODO: Проверить, что метод GET (иначе 405)
// --TODO: Получить заголовок Authorization из r.Header
// --TODO: Проверить, что он начинается с "Bearer "
// --TODO: Извлечь токен — всё, что после "Bearer "
// --TODO: Вызвать h.AuthService.GetUserFromToken(token) — получить пользователя
// --TODO: Если токен невалиден — вернуть 401
// --TODO: Если токен валиден — получить доски через h.BoardService.GetBoards()
// --TODO: Вернуть 200 с JSON: { "user": { "id", "email", "username" }, "boards": [...] }
func (bh *BoardHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetBoard")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		log.Println("Cookie not found")
		return
	}

	tokenString := cookie.Value
	user, err := bh.AuthService.GetUserFromToken(ctx, tokenString)
	//User, err := h.AuthService.GetCurrentUser()
	if err != nil {
		http.Error(w, "401 : "+err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}
	boards, err := bh.BoardService.GetBoards(ctx, user.ID)
	if err != nil {
		log.Printf("error while get boards: %s", err)
		return
	}
	json_Boards, errB := json.Marshal(boards)
	if errB != nil {
		log.Printf("error while serialize Boards: %s", errB)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_Boards)

}
func (bh *BoardHandler) CreateBoards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	newBoard := new(models.Board)
	err := decoder.Decode(newBoard)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	err = bh.BoardService.AddBoard(ctx, newBoard)
	if err != nil {
		log.Printf("error while create Board: %s", err)
		http.Error(w, "400 :"+err.Error(), http.StatusBadRequest)
		return
	}
}

func (bh *BoardHandler) BoardDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	err := bh.BoardService.DeleteBoard(ctx, vars)
	if err != nil {
		log.Printf("error while delete Board: %s", err)
		http.Error(w, "400 :"+err.Error(), http.StatusBadRequest)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (bh *BoardHandler) BoardRestore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	err := bh.BoardService.RestoreBoard(ctx, vars)
	if err != nil {
		log.Printf("error while restore Board: %s", err)
		http.Error(w, "400 :"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (bh *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	boardID, err := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	if err != nil {
		http.Error(w, "invalid board_id", http.StatusBadRequest)
		return
	}

	// Получаем полные данные доски
	fullBoard, err := bh.BoardService.GetBoard(ctx, boardID)
	if err != nil {
		http.Error(w, "failed to get board data", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fullBoard); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
func (bh *BoardHandler) RenameBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	decoder := json.NewDecoder(r.Body)
	newBoard := new(models.Board)
	err := decoder.Decode(newBoard)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	newBoard.ID = vars
	log.Println(newBoard)
	board, err := bh.BoardService.RenameBoard(ctx, newBoard)
	if err != nil {
		http.Error(w, "failed to rename board", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(board); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
func (bh *BoardHandler) ArchivedBoard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	err := bh.BoardService.ArchivedBoard(ctx, vars)
	if err != nil {
		http.Error(w, "failed to archived board", http.StatusInternalServerError)
		return
	}
}

func (bh *BoardHandler) Board(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bh.GetBoard(w, r)
	case http.MethodDelete:
		bh.BoardDelete(w, r)
	case http.MethodPut:
		bh.RenameBoard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (bh *BoardHandler) CreateOrGetBoards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		bh.CreateBoards(w, r)
	case http.MethodGet:
		bh.GetBoards(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Запрос " + r.Method + ",а должен быть GET или POST")
		return
	}
}
