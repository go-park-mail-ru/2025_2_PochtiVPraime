package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

type ChecklistPointHandler struct {
	ChecklistPointService services.ChecklistPointService
	AuthService           services.AuthService
}

func NewChecklistPointHandler(checklistPointService *services.ChecklistPointService, authService *services.AuthService) *ChecklistPointHandler {
	return &ChecklistPointHandler{
		ChecklistPointService: *checklistPointService,
		AuthService:           *authService,
	}
}

// GetChecklistPoints возвращает все пункты чеклиста
func (cph *ChecklistPointHandler) GetChecklistPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetChecklistPoints")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}

	points, err := cph.ChecklistPointService.GetChecklistPointsByChecklistID(ctx, checklistID)
	if err != nil {
		log.Printf("error while get checklist points: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonPoints, err := json.Marshal(points)
	if err != nil {
		log.Printf("error while serialize checklist points: %s", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPoints)
}

// CreateChecklistPoint создает новый пункт чеклиста
func (cph *ChecklistPointHandler) CreateChecklistPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("CreateChecklistPoint")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	newPoint := new(models.ChecklistPoint)
	err = decoder.Decode(newPoint)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}
	err = cph.ChecklistPointService.CreateChecklistPoint(ctx, newPoint, checklistID)
	if err != nil {
		log.Printf("error while create checklist point: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPoint)
}

// GetChecklistPoint возвращает пункт чеклиста по ID
func (cph *ChecklistPointHandler) GetChecklistPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("GetChecklistPoint")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	pointID, err := strconv.ParseInt(r.PathValue("pointId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID пункта чеклиста", http.StatusBadRequest)
		return
	}

	point, err := cph.ChecklistPointService.GetChecklistPointByID(ctx, pointID)
	if err != nil {
		log.Printf("error while get checklist point: %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(point)
}

// UpdateChecklistPoint обновляет пункт чеклиста
func (cph *ChecklistPointHandler) UpdateChecklistPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("UpdateChecklistPoint")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	pointID, err := strconv.ParseInt(r.PathValue("pointId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID пункта чеклиста", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	updatedPoint := new(models.ChecklistPoint)
	err = decoder.Decode(updatedPoint)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	updatedPoint.ID = pointID

	err = cph.ChecklistPointService.UpdateChecklistPoint(ctx, updatedPoint)
	if err != nil {
		log.Printf("error while update checklist point: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPoint)
}

// UpdateChecklistPointStatus обновляет статус checked пункта чеклиста
func (cph *ChecklistPointHandler) UpdateChecklistPointStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("UpdateChecklistPointStatus")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	pointID, err := strconv.ParseInt(r.PathValue("pointId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID пункта чеклиста", http.StatusBadRequest)
		return
	}

	var request struct {
		Checked bool `json:"checked"`
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err = cph.ChecklistPointService.UpdateCheckedStatus(ctx, pointID, request.Checked)
	if err != nil {
		log.Printf("error while update checklist point status: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "checked status updated"})
}

// DeleteChecklistPoint удаляет пункт чеклиста
func (cph *ChecklistPointHandler) DeleteChecklistPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("DeleteChecklistPoint")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	pointID, err := strconv.ParseInt(r.PathValue("pointId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID пункта чеклиста", http.StatusBadRequest)
		return
	}

	err = cph.ChecklistPointService.DeleteChecklistPoint(ctx, pointID)
	if err != nil {
		log.Printf("error while delete checklist point: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// DeleteChecklistPoints удаляет все пункты чеклиста
func (cph *ChecklistPointHandler) DeleteChecklistPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("DeleteChecklistPoints")

	// Получаем пользователя из токена
	_, err := cph.getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}

	checklistID, err := strconv.ParseInt(r.PathValue("checklistId"), 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID чеклиста", http.StatusBadRequest)
		return
	}

	err = cph.ChecklistPointService.DeletePointsByChecklistID(ctx, checklistID)
	if err != nil {
		log.Printf("error while delete checklist points: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "all points deleted"})
}

// ChecklistPoint обрабатывает различные HTTP методы для работы с пунктом чеклиста
func (cph *ChecklistPointHandler) ChecklistPoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cph.GetChecklistPoint(w, r)
	case http.MethodPut:
		cph.UpdateChecklistPoint(w, r)
	case http.MethodDelete:
		cph.DeleteChecklistPoint(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (cph *ChecklistPointHandler) GetOrCreateChecklistsPoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cph.GetChecklistPoints(w, r)
	case http.MethodPost:
		cph.CreateChecklistPoint(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// ChecklistPoints обрабатывает GET запрос для получения всех пунктов чеклиста
func (cph *ChecklistPointHandler) ChecklistPoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cph.GetChecklistPoints(w, r)
	case http.MethodDelete:
		cph.DeleteChecklistPoints(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// ChecklistPointStatus обрабатывает PATCH запрос для обновления статуса пункта
func (cph *ChecklistPointHandler) ChecklistPointStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cph.UpdateChecklistPointStatus(w, r)
}

// getUserFromRequest вспомогательный метод для получения пользователя из запроса
func (cph *ChecklistPointHandler) getUserFromRequest(r *http.Request) (*models.User, error) {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	tokenString := cookie.Value
	user, err := cph.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return user, nil
}
