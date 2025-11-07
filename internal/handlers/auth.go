package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

// AuthHandler — обработчик HTTP-запросов аунтификации
// Внедряет зависимости: AuthService
type AuthHandler struct {
	AuthService services.AuthService
	JWT         string
}

// NewAuthHandler — конструктор для Dependency Injection
// Создаёт и инициализирует все зависимости
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: *authService,
	}
}

// Register — обрабатывает POST api/auth/register
// --TODO: Проверить, что метод POST (иначе 405)
// --TODO: Декодировать JSON из тела запроса (email, username, password)
// --TODO: Проверить, что email не пустой и содержит "@"
// --TODO: Проверить, что username не пустой и не слишком длинный
// --TODO: Проверить, что password не короче 6 символов
// TODO: Проверить, что пользователь с таким email уже не существует
// --TODO: Вызвать h.AuthService.Register(email, username, password)
// --TODO: Если ошибка — вернуть 400 с сообщением об ошибке
// TODO: Если успех — вернуть 201 с JSON: { "user": { "id", "email", "username" } }
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Пока просто отвечаем заглушкой
	ctx := r.Context()
	if r.Method != http.MethodPost {
		log.Printf("Запрос " + r.Method + ",а должен быть POST")
		http.Error(w, "405 : NotAcceptable", http.StatusNotAcceptable)
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	newUserInput := new(models.User)
	err := decoder.Decode(newUserInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	validUser, err := h.AuthService.Register(ctx, newUserInput)
	if err != nil {
		log.Printf("error while register User in service: %s", err)
		http.Error(w, "400 : "+err.Error(), http.StatusBadRequest)
		return
	}
	json_User, err := json.Marshal(validUser)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
	//log.Printf(string(json_User))
	w.Header().Set("Content-Type", "application/json")
}

// Login — обрабатывает POST api/auth/login
// --TODO: Проверить, что метод POST (иначе 405)
// --TODO: Декодировать JSON из тела запроса (email, password)
// --TODO: Проверить, что email и password не пустые
// --TODO: Вызвать h.AuthService.Login(email, password) — получить JWT токен
// --TODO: Если ошибка — вернуть 401 с сообщением "неправильный email или пароль"
// --TODO: Если успех — вернуть 200 с JSON: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		log.Printf("Запрос " + r.Method + ",а должен быть POST")
		http.Error(w, "405 : NotAcceptable", http.StatusNotAcceptable)
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	UserInput := new(models.User)
	err := decoder.Decode(UserInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	log.Println(UserInput)
	JWT, err := h.AuthService.Login(ctx, UserInput)
	h.JWT = JWT
	if err != nil {
		log.Printf("error while authorizate: %s", err)
		http.Error(w, "401 : "+err.Error(), http.StatusUnauthorized)
		return
	}
	//установка куки
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    JWT, //сюда записать токен
		Path:     "/",
		HttpOnly: true,                    // Доступ только через HTTP, защита от XSS
		Secure:   false,                   // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
		MaxAge:   900,                     // время жизни куки в секундах (поставил 15 минут)
	}
	http.SetCookie(w, cookie)
	log.Println("cookie created")
	w.Header().Set("Content-Type", "application/json")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Удаление куки
	}
	http.SetCookie(w, cookie)
	h.AuthService.Logout()
	log.Println("cookie deleted")
	w.Write([]byte("Cookie deleted!"))
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}
	tokenString := cookie.Value
	user, err := h.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		http.Error(w, "401 : "+err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}

	//user.Password = []byte{}
	json_User, err := json.Marshal(user)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
	//log.Printf(string(json_User))
}

func (h *AuthHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPut {
		log.Printf("Запрос " + r.Method + ",а должен быть POST")
		http.Error(w, "405 : NotAcceptable", http.StatusNotAcceptable)
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	newUserInput := new(models.User)
	err := decoder.Decode(newUserInput)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}
	tokenString := cookie.Value
	user, err := h.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		http.Error(w, "401 : "+err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}
	newUserInput.ID = user.ID

	user, err = h.AuthService.UpdateUser(ctx, newUserInput)
	if err != nil {
		log.Printf("error while Update User in service: %s", err)
		http.Error(w, "400 : "+err.Error(), http.StatusBadRequest)
		return
	}
	json_User, err := json.Marshal(user)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
}

func (h *AuthHandler) PasswordUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPut {
		log.Printf("Запрос " + r.Method + ",а должен быть POST")
		http.Error(w, "405 : NotAcceptable", http.StatusNotAcceptable)
		return
	}

	type ChangePass struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	passwords := new(ChangePass)
	err := decoder.Decode(passwords)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}
	tokenString := cookie.Value
	user, err := h.AuthService.GetUserFromToken(ctx, tokenString)
	if err != nil {
		http.Error(w, "401 : "+err.Error(), http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}

	user, err = h.AuthService.PasswordUpdate(ctx, passwords.OldPassword, passwords.NewPassword, user.ID)
	if err != nil {
		log.Printf("error while Update User in service: %s", err)
		http.Error(w, "400 : "+err.Error(), http.StatusBadRequest)
		return
	}
	json_User, err := json.Marshal(user)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
}
