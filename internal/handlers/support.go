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
// Внедряет зависимости: SupportService
type SupportHandler struct {
	SupportService services.SupportService
	AuthService    services.AuthService
}

// Создаёт и инициализирует все зависимости
func NewSupportHandler(supportService *services.SupportService, authService *services.AuthService) *SupportHandler {
	return &SupportHandler{
		SupportService: *supportService,
		AuthService:    *authService,
	}
}

// GetBoards — обрабатывает GET /api/boards
func (sh *SupportHandler) GetUserSupportForms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//log.Println("GetBoard")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		log.Println("Cookie not found")
		return
	}

	tokenString := cookie.Value
	user, err := sh.AuthService.GetUserFromToken(ctx, tokenString)
	//User, err := h.AuthService.GetCurrentUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}
	forms, err := sh.SupportService.GetUserSupportForms(ctx, user.ID)
	if err != nil {
		log.Printf("error while get forms: %s", err)
		return
	}
	json_Boards, errB := json.Marshal(forms)
	if errB != nil {
		log.Printf("error while serialize Boards: %s", errB)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_Boards)

}

func (sh *SupportHandler) GetAllSupportForms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//log.Println("GetBoard")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		log.Println("Cookie not found")
		return
	}

	tokenString := cookie.Value
	_, err = sh.AuthService.GetUserFromToken(ctx, tokenString)
	//User, err := h.AuthService.GetCurrentUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}
	forms, err := sh.SupportService.GetAllSupportForms(ctx)
	if err != nil {
		log.Printf("error while get forms: %s", err)
		return
	}
	json_Boards, errB := json.Marshal(forms)
	if errB != nil {
		log.Printf("error while serialize Boards: %s", errB)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_Boards)

}

func (sh *SupportHandler) CreateSupportForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	newSupportForm := new(models.SupportForm)
	err := decoder.Decode(newSupportForm)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		log.Println("Cookie not found")
		return
	}

	tokenString := cookie.Value
	user, err := sh.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}
	newSupportForm.UserId = user.ID
	err = sh.SupportService.AddSupportForm(ctx, newSupportForm)
	if err != nil {
		log.Printf("ошибка при создании формы: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (sh *SupportHandler) DeleteSupportForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars, _ := strconv.ParseInt(r.PathValue("formId"), 10, 64)
	supportForm, err := sh.SupportService.GetSupportFormById(ctx, vars)
	if err != nil {
		log.Printf("error while get support form: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		log.Println("Cookie not found")
		return
	}

	tokenString := cookie.Value
	user, err := sh.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization error:", err)
		return
	}
	if supportForm.UserId != user.ID {
		http.Error(w, "пользователь не является хозяином формы", http.StatusBadRequest)
		return
	}
	err = sh.SupportService.DeleteSupportForm(ctx, vars)
	if err != nil {
		log.Printf("error while delete support form: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (sh *SupportHandler) GetSupportForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	supportFormId, err := strconv.ParseInt(r.PathValue("formId"), 10, 64)
	if err != nil {
		http.Error(w, "invalid board_id", http.StatusBadRequest)
		return
	}

	// Получаем полные данные доски
	supportForm, err := sh.SupportService.GetSupportFormById(ctx, supportFormId)
	if err != nil {
		http.Error(w, "failed to get board data", http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(supportForm); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (sh *SupportHandler) CreateOrGetForms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		sh.CreateSupportForm(w, r)
	case http.MethodGet:
		sh.GetUserSupportForms(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Запрос " + r.Method + ",а должен быть GET или POST")
		return
	}
}

func (sh *SupportHandler) DeleteOrGetForm(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		sh.DeleteSupportForm(w, r)
	case http.MethodGet:
		sh.GetSupportForm(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Запрос " + r.Method + ",а должен быть GET или POST")
		return
	}
}
