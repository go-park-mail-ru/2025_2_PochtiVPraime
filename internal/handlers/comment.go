package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

type CommentHandler struct {
	CommentService services.CommentService
	AuthService    services.AuthService
}

func NewCommentHandler(commentService *services.CommentService, authService *services.AuthService) *CommentHandler {
	return &CommentHandler{
		CommentService: *commentService,
		AuthService:    *authService,
	}
}

// GetComments возвращает все комментарии карточки
func (ch *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetComments")

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID карточки", http.StatusBadRequest)
		return
	}

	// Здесь нужно получить boardMemberID для пользователя и карточки
	// Для упрощения используем user.ID как boardMemberID
	// В реальном приложении нужно получить правильный boardMemberID

	comments, err := ch.CommentService.GetCommentsByCardID(ctx, cardID)
	if err != nil {
		log.Printf("error while get comments: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonComments, err := json.Marshal(comments)
	if err != nil {
		log.Printf("error while serialize comments: %s", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonComments)
}

// CreateComment создает новый комментарий
func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("CreateComment")

	// Получаем пользователя из токена
	user, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID карточки", http.StatusBadRequest)
		return
	}

	var request struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Создаем комментарий
	comment := &models.Comment{
		CardId:  cardID,
		Content: request.Content,
	}

	// Используем user.ID как boardMemberID (в реальном приложении нужно получить правильный boardMemberID)
	err = ch.CommentService.CreateComment(ctx, comment, user.ID)
	if err != nil {
		log.Printf("error while create comment: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetComment возвращает комментарий по ID
func (ch *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetComment")

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	commentID, err := strconv.ParseInt(r.PathValue("commentId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := ch.CommentService.GetCommentByID(ctx, commentID)
	if err != nil {
		log.Printf("error while get comment: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// UpdateComment обновляет комментарий
func (ch *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("UpdateComment")

	// Получаем пользователя из токена
	user, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	commentID, err := strconv.ParseInt(r.PathValue("commentId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	var request struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	canEdit, err := ch.CommentService.CanUserEditComment(ctx, commentID, user.ID)
	if err != nil {
		log.Printf("error while check edit permissions: %s", err)
		http.Error(w, "Ошибка при проверке прав доступа", http.StatusInternalServerError)
		return
	}
	if !canEdit {
		http.Error(w, "Недостаточно прав для редактирования комментария", http.StatusForbidden)
		return
	}

	// Получаем текущий комментарий
	comment, err := ch.CommentService.GetCommentByID(ctx, commentID)
	if err != nil {
		log.Printf("error while get comment: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Обновляем контент
	comment.Content = request.Content
	err = ch.CommentService.UpdateComment(ctx, comment)
	if err != nil {
		log.Printf("error while update comment: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// DeleteComment удаляет комментарий
func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("DeleteComment")

	// Получаем пользователя из токена
	user, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	commentID, err := strconv.ParseInt(r.PathValue("commentId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	canDelete, err := ch.CommentService.CanUserDeleteComment(ctx, commentID, user.ID)
	if err != nil {
		log.Printf("error while check delete permissions: %s", err)
		http.Error(w, "Ошибка при проверке прав доступа", http.StatusInternalServerError)
		return
	}
	if !canDelete {
		http.Error(w, "Недостаточно прав для удаления комментария", http.StatusForbidden)
		return
	}

	err = ch.CommentService.DeleteComment(ctx, commentID)
	if err != nil {
		log.Printf("error while delete comment: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// UpdateCommentContent обновляет только содержание комментария
func (ch *CommentHandler) UpdateCommentContent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("UpdateCommentContent")

	// Получаем пользователя из токена
	user, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	commentID, err := strconv.ParseInt(r.PathValue("commentId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	var request struct {
		Content string `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	canEdit, err := ch.CommentService.CanUserEditComment(ctx, commentID, user.ID)
	if err != nil {
		log.Printf("error while check edit permissions: %s", err)
		http.Error(w, "Ошибка при проверке прав доступа", http.StatusInternalServerError)
		return
	}
	if !canEdit {
		http.Error(w, "Недостаточно прав для редактирования комментария", http.StatusForbidden)
		return
	}

	err = ch.CommentService.UpdateCommentContent(ctx, commentID, request.Content)
	if err != nil {
		log.Printf("error while update comment content: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "content updated"})
}

// GetUserComments возвращает все комментарии пользователя
func (ch *CommentHandler) GetUserComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetUserComments")

	// Получаем пользователя из токена
	user, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	comments, err := ch.CommentService.GetCommentsByBoardMemberID(ctx, user.ID)
	if err != nil {
		log.Printf("error while get user comments: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonComments, err := json.Marshal(comments)
	if err != nil {
		log.Printf("error while serialize user comments: %s", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonComments)
}

// Comment обрабатывает различные HTTP методы для работы с комментарием
func (ch *CommentHandler) Comment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetComment(w, r)
	case http.MethodPut:
		ch.UpdateComment(w, r)
	case http.MethodDelete:
		ch.DeleteComment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Comments обрабатывает GET запрос для получения всех комментариев карточки
func (ch *CommentHandler) Comments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetComments(w, r)
	case http.MethodPost:
		ch.CreateComment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// CommentContent обрабатывает PATCH запрос для обновления содержания комментария
func (ch *CommentHandler) CommentContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ch.UpdateCommentContent(w, r)
}

// UserComments обрабатывает GET запрос для получения комментариев пользователя
func (ch *CommentHandler) UserComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ch.GetUserComments(w, r)
}

// getUserFromRequest вспомогательный метод для получения пользователя из запроса
func (ch *CommentHandler) getUserFromRequest(r *http.Request) (*models.User, error) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	user, err := ch.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return user, nil
}
