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
func NewListHandler(listService services.ListService, authService services.AuthService) *ListHandler {
	return &ListHandler{
		ListService: listService,
		AuthService: authService,
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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из токена
	user, err := lh.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	// Парсим тело запроса
	var requestData struct {
		Title   string `json:"title"`
		BoardID int64  `json:"boardId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %s", err)
		return
	}

	// Валидация данных
	if requestData.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if requestData.BoardID == 0 {
		http.Error(w, "BoardID is required", http.StatusBadRequest)
		return
	}

	// Создаем список через сервис
	list, err := lh.ListService.CreateList(ctx, requestData.Title, requestData.BoardID, user.ID)
	if err != nil {
		log.Printf("Error creating list: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный список
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

// GetLists обрабатывает GET /api/boards/{boardId}/lists
func (lh *ListHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем boardId из query параметров
	boardIDStr := r.URL.Query().Get("boardId")
	if boardIDStr == "" {
		// Пробуем извлечь из path (если используется роутинг типа /boards/{id}/lists)
		pathParts := strings.Split(r.URL.Path, "/")
		for i, part := range pathParts {
			if part == "boards" && i+1 < len(pathParts) {
				boardIDStr = pathParts[i+1]
				break
			}
		}
	}

	if boardIDStr == "" {
		http.Error(w, "BoardID is required", http.StatusBadRequest)
		return
	}

	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
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

// UpdateList обрабатывает PUT /api/lists/{listId}
func (lh *ListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из токена
	user, err := lh.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	// Извлекаем listId из URL
	pathParts := strings.Split(r.URL.Path, "/")
	var listIDStr string
	for i, part := range pathParts {
		if part == "lists" && i+1 < len(pathParts) {
			listIDStr = pathParts[i+1]
			break
		}
	}

	if listIDStr == "" {
		http.Error(w, "ListID is required", http.StatusBadRequest)
		return
	}

	listID, err := strconv.ParseInt(listIDStr, 10, 64)
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

	// Обновляем список через сервис
	updatedList, err := lh.ListService.UpdateList(ctx, &list, user.ID)
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

// DeleteList обрабатывает DELETE /api/lists/{listId}
func (lh *ListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из токена
	user, err := lh.GetUserFromRequest(ctx, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	// Извлекаем listId из URL
	pathParts := strings.Split(r.URL.Path, "/")
	var listIDStr string
	for i, part := range pathParts {
		if part == "lists" && i+1 < len(pathParts) {
			listIDStr = pathParts[i+1]
			break
		}
	}

	if listIDStr == "" {
		http.Error(w, "ListID is required", http.StatusBadRequest)
		return
	}

	listID, err := strconv.ParseInt(listIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ListID", http.StatusBadRequest)
		return
	}

	// Удаляем список через сервис
	err = lh.ListService.DeleteListWithCard(ctx, listID, user.ID)
	if err != nil {
		log.Printf("Error deleting list: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusNoContent)
}
