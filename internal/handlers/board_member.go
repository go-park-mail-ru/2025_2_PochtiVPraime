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
type BoardMemberHandler struct {
	BoardMemberService services.BoardMemberService
}

// NewHandler — конструктор для Dependency Injection
// Создаёт и инициализирует все зависимости
func NewBoardMemberHandler(boardMemberService *services.BoardMemberService) *BoardMemberHandler {
	return &BoardMemberHandler{
		BoardMemberService: *boardMemberService,
	}
}

func (bmh *BoardMemberHandler) CreateBoardMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("CreateBoardMember")
	decoder := json.NewDecoder(r.Body)
	newBoardMember := new(models.BoardMember)
	err := decoder.Decode(newBoardMember)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	zaglushka := int64(1)
	boardId, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	newBoardMember, err = bmh.BoardMemberService.CreateBoardMember(ctx, newBoardMember, zaglushka, boardId)
	if err != nil {
		log.Printf("ошибка при создании роли участника доски: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (bmh *BoardMemberHandler) GetBoardMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetBoardMembers")
	vars, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	boardMembers, err := bmh.BoardMemberService.GetBoardMembers(ctx, vars)
	if err != nil {
		log.Printf("ошибка при получении участников доски: %s", err)
		return
	}
	json_BoardMembers, err := json.Marshal(boardMembers)
	if err != nil {
		log.Printf("ошибка во время сериализации участников доски s: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_BoardMembers)

}

/*
	func (bmh *BoardMemberHandler) GetBoardMember(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("GetBoardMembers")
		vars, _ := strconv.ParseInt(r.PathValue("bordMembersId"), 10, 64)
		boardMembers, err := bmh.BoardMemberService.GetBoardMemberByUserId(ctx, vars)
		if err != nil {
			log.Printf("ошибка при получении участников доски: %s", err)
			return
		}
		json_BoardMembers, err := json.Marshal(boardMembers)
		if err != nil {
			log.Printf("ошибка во время сериализации участников доски s: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_BoardMembers)

}
*/
func (bmh *BoardMemberHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("ChangeRole")
	decoder := json.NewDecoder(r.Body)
	newBoardMember := new(models.BoardMember)
	err := decoder.Decode(newBoardMember)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	zaglushka := int64(1)
	boardId, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	err = bmh.BoardMemberService.ChangeRole(ctx, zaglushka, boardId, newBoardMember)
	if err != nil {
		log.Printf("ошибка при получении изменении роли участника доски: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (bmh *BoardMemberHandler) DeleteBoardMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("DeleteBoardMember")
	decoder := json.NewDecoder(r.Body)
	newBoardMember := new(models.BoardMember)
	err := decoder.Decode(newBoardMember)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	zaglushka := int64(1)
	boardId, _ := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	err = bmh.BoardMemberService.DeleteBoardMember(ctx, zaglushka, boardId, newBoardMember)
	if err != nil {
		log.Printf("ошибка при получении изменении роли участника доски: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (bmh *BoardMemberHandler) BoardMember(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		bmh.CreateBoardMember(w, r)
	case http.MethodGet:
		bmh.GetBoardMembers(w, r)
	case http.MethodDelete:
		bmh.DeleteBoardMember(w, r)
	case http.MethodPut:
		bmh.ChangeRole(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Запрос " + r.Method + ",а должен быть GET или POST")
		return
	}
}
