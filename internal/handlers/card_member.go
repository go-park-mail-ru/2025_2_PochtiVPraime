package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

type CardMemberHandler struct {
	CardMemberService services.CardMemberService
	AuthService       services.AuthService
}

func NewCardMemberHandler(cardMemberService *services.CardMemberService, authService *services.AuthService) *CardMemberHandler {
	return &CardMemberHandler{
		CardMemberService: *cardMemberService,
		AuthService:       *authService,
	}
}

// GetCardMembers возвращает всех участников карточки
func (cmh *CardMemberHandler) GetCardMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetCardMembers")

	// Получаем пользователя из токена
	_, err := cmh.getUserFromRequest(r)
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

	cardMembers, err := cmh.CardMemberService.GetCardMembersByCardID(ctx, cardID)
	if err != nil {
		log.Printf("error while get card members: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonCardMembers, err := json.Marshal(cardMembers)
	if err != nil {
		log.Printf("error while serialize card members: %s", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCardMembers)
}

// CreateCardMember создает связь между карточкой и участником
func (cmh *CardMemberHandler) CreateCardMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("CreateCardMember")

	// Получаем пользователя из токена
	_, err := cmh.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	newCardMember := new(models.CardMember)
	err = decoder.Decode(newCardMember)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	cardId, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID связи карточки и участника", http.StatusBadRequest)
		return
	}

	err = cmh.CardMemberService.CreateCardMember(ctx, newCardMember, cardId)
	if err != nil {
		log.Printf("error while create card member: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCardMember)
}

// GetCardMember возвращает связь по ID
func (cmh *CardMemberHandler) GetCardMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetCardMember")

	// Получаем пользователя из токена
	_, err := cmh.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	cardMemberID, err := strconv.ParseInt(r.PathValue("cardMemberId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID связи карточки и участника", http.StatusBadRequest)
		return
	}

	cardMember, err := cmh.CardMemberService.GetCardMemberByID(ctx, cardMemberID)
	if err != nil {
		log.Printf("error while get card member: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cardMember)
}

// DeleteCardMember удаляет связь между карточкой и участником
func (cmh *CardMemberHandler) DeleteCardMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("DeleteCardMember")

	// Получаем пользователя из токена
	_, err := cmh.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID связи карточки и участника", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	cardMember := new(models.CardMember)
	err = decoder.Decode(cardMember)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err = cmh.CardMemberService.DeleteCardMember(ctx, cardID, cardMember.BoardMemberID)
	if err != nil {
		log.Printf("error while delete card member: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// CardMember обрабатывает различные HTTP методы для работы со связью карточки и участника
func (cmh *CardMemberHandler) CardMember(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cmh.CreateCardMember(w, r)
	case http.MethodGet:
		cmh.GetCardMembers(w, r)
	case http.MethodDelete:
		cmh.DeleteCardMember(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

/*
// CardMembers обрабатывает GET запрос для получения участников карточки
func (cmh *CardMemberHandler) DeleteOrGetCardMember(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		cmh.DeleteCardMember(w, r)
	case http.MethodGet:
		cmh.GetCardMember(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
*/
// getUserFromRequest вспомогательный метод для получения пользователя из запроса
func (cmh *CardMemberHandler) getUserFromRequest(r *http.Request) (*models.User, error) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	user, err := cmh.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return user, nil
}
