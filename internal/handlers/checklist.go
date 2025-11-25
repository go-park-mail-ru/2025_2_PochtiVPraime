package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

type ChecklistHandler struct {
	ChecklistService services.ChecklistService
	AuthService      services.AuthService
}

func NewChecklistHandler(checklistService *services.ChecklistService, authService *services.AuthService) *ChecklistHandler {
	return &ChecklistHandler{
		ChecklistService: *checklistService,
		AuthService:      *authService,
	}
}

// CreateChecklist создает новый чеклист
func (ch *ChecklistHandler) CreateChecklist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	/*
		// Получаем пользователя из токена
		user, err := ch.getUserFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	*/
	decoder := json.NewDecoder(r.Body)
	newChecklist := new(models.Checklist)
	err := decoder.Decode(newChecklist)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Здесь можно добавить проверку прав доступа пользователя к карточке
	// если нужно убедиться, что пользователь имеет доступ к карточке, к которой добавляется чеклист

	err = ch.ChecklistService.CreateChecklist(ctx, newChecklist)
	if err != nil {
		log.Printf("error while create checklist: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newChecklist)
}

// GetChecklist возвращает чеклист по ID
func (ch *ChecklistHandler) GetChecklist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	checklistId, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}
	checklist, err := ch.ChecklistService.GetChecklistByID(ctx, checklistId)
	if err != nil {
		log.Printf("error while get checklist: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checklist)
}

// GetChecklistByCard возвращает чеклист по ID карточки
func (ch *ChecklistHandler) GetChecklistsByCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cardID, err := strconv.ParseInt(r.PathValue("taskId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID карточки", http.StatusBadRequest)
		return
	}

	checklists, err := ch.ChecklistService.GetChecklistsByCardID(ctx, cardID)
	if err != nil {
		log.Printf("error while get checklist by card: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checklists)
}

// UpdateChecklist обновляет чеклист
func (ch *ChecklistHandler) UpdateChecklist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	updatedChecklist := new(models.Checklist)
	err = decoder.Decode(updatedChecklist)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err = ch.ChecklistService.UpdateChecklist(ctx, updatedChecklist, checklistID)
	if err != nil {
		log.Printf("ошибка во время обновления чеклиста: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedChecklist)
}

// DeleteChecklist удаляет чеклист
func (ch *ChecklistHandler) DeleteChecklist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Получаем пользователя из токена
	_, err := ch.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}

	err = ch.ChecklistService.DeleteChecklist(ctx, checklistID)
	if err != nil {
		log.Printf("error while delete checklist: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Checklist обрабатывает различные HTTP методы для работы с чеклистом
func (ch *ChecklistHandler) GetOrCreateChecklists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ch.CreateChecklist(w, r)
	case http.MethodGet:
		ch.GetChecklistsByCard(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Checklist обрабатывает различные HTTP методы для работы с чеклистом
func (ch *ChecklistHandler) Checklist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ch.GetChecklist(w, r)
	case http.MethodPut:
		ch.UpdateChecklist(w, r)
	case http.MethodDelete:
		ch.DeleteChecklist(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// getUserFromRequest вспомогательный метод для получения пользователя из запроса
func (ch *ChecklistHandler) getUserFromRequest(r *http.Request) (*models.User, error) {
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
