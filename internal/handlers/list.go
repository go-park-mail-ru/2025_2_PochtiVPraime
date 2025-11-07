package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

// ListHandler — обработчик HTTP-запросов для списков
type ListHandler struct {
	ListService services.ListService
	AuthService services.AuthService
}

// NewListHandler — конструктор для Dependency Injection
func NewListHandler(listService *services.ListService, authService *services.AuthService) *ListHandler {
	return &ListHandler{
		ListService: *listService,
		AuthService: *authService,
	}
}

// GetUserFromRequest извлекает пользователя из токена
func (lh *ListHandler) GetUserFromRequest(ctx context.Context, r *http.Request) (*models.User, error) {
	// Получаем токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, http.ErrNoCookie
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.ErrNoCookie
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	return lh.AuthService.GetUserFromToken(ctx, tokenString)
}

// CreateList обрабатывает POST /api/lists
func (lh *ListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	/*
		// Получаем пользователя из токена
		user, err := lh.GetUserFromRequest(ctx, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Authorization error:", err)
			return
		}
	*/
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %s", err)
		return
	}

	// Извлекаем boardId из query параметров
	boardID, err := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid BoardID", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if list.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	list.BoardId = boardID
	zaglushka := 12 //поменять, когда разберусь почему не видет куки
	// Создаем список через сервис
	newList, err := lh.ListService.CreateList(ctx, list.Title, list.BoardId, int64(zaglushka))
	if err != nil {
		log.Printf("Error creating list: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный список
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newList); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

// GetLists обрабатывает GET /api/boards/{boardId}/lists
func (lh *ListHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Извлекаем boardId из query параметров
	boardID, err := strconv.ParseInt(r.PathValue("boardId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid BoardID", http.StatusBadRequest)
		return
	}

	// Получаем списки через сервис
	lists, err := lh.ListService.GetLists(ctx, boardID)
	if err != nil {
		log.Printf("Error getting lists: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем списки
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lists); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

// UpdateList обрабатывает PUT /board/{boardId}/lists/{listId}
func (lh *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	/*
		// Получаем пользователя из токена
		user, err := lh.GetUserFromRequest(ctx, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Authorization error:", err)
			return
		}
	*/

	listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ListID", http.StatusBadRequest)
		return
	}

	// Парсим тело запроса
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %s", err)
		return
	}

	// Устанавливаем ID из URL
	list.ID = listID

	// Валидация данных
	if list.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	zaglushka := 1
	// Обновляем список через сервис
	updatedList, err := lh.ListService.UpdateList(ctx, &list, int64(zaglushka))
	if err != nil {
		log.Printf("Error updating list: %s", err)
		// Можно добавить проверку на конкретные ошибки (например, "not found", "access denied")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем обновленный список
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedList); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

// GetList получаем список по Id GET /api/lists/{listId}
func (lh *ListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Извлекаем boardId из query параметров
	listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid listID", http.StatusBadRequest)
		return
	}
	list, err := lh.ListService.GetList(ctx, listID)
	if err != nil {
		http.Error(w, "failed to get list", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

// DeleteList обрабатывает DELETE /api/lists/{listId}
func (lh *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	/*
		// Получаем пользователя из токена
		user, err := lh.GetUserFromRequest(ctx, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("Authorization error:", err)
			return
		}
	*/

	listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ListID", http.StatusBadRequest)
		return
	}

	zaglushka := 1
	// Удаляем список через сервис
	err = lh.ListService.DeleteListWithCard(ctx, listID, int64(zaglushka))
	if err != nil {
		log.Printf("Error deleting list: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusNoContent)
}

func (lh *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		lh.GetList(w, r)
	case http.MethodDelete:
		lh.DeleteList(w, r)
	case http.MethodPut:
		lh.UpdateList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (lh *ListHandler) CreateOrGetLists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		lh.GetLists(w, r)
	case http.MethodPost:
		lh.CreateList(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
