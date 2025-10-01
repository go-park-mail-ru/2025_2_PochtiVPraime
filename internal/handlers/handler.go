package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/services"
)

// Handler — обработчик HTTP-запросов
// Внедряет зависимости: BoardService и AuthService
type Handler struct {
	BoardService *services.BoardService
	AuthService  *services.AuthService
}

// NewHandler — конструктор для Dependency Injection
// Создаёт и инициализирует все зависимости
func NewHandler() *Handler {
	return &Handler{
		BoardService: services.NewBoardService(),
		AuthService:  services.NewAuthService(),
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
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Пока просто отвечаем заглушкой
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

	email := newUserInput.Email
	if !strings.Contains(email, "@") || len(email) == 0 { //наверное len должна быть хотябы 6
		newErr := errors.New("Не содержит @ или слишком короткий email") //тк len(a@b.ru)
		log.Printf("error while email not valid: %s", newErr)
		return
	}
	username := newUserInput.Username
	if len(username) <= 0 || len(username) > 25 {
		newErr := errors.New("слишком короткое или слишком длинное имя")
		log.Printf("error while name not valid: %s", newErr)
		return
	}

	password := newUserInput.Password
	if len(password) < 6 {
		newErr := errors.New("слишком короткий пароль")
		log.Printf("error while name not valid: %s", newErr)
		return
	}

	validUser, err := h.AuthService.Register(email, username, password)
	if err != nil {
		log.Printf("error while saving User in Service: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	validUser.Password = ""
	json_User, err := json.Marshal(validUser)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
	log.Printf(string(json_User))
	w.Header().Set("Content-Type", "application/json")
}

// Login — обрабатывает POST api/auth/login
// --TODO: Проверить, что метод POST (иначе 405)
// --TODO: Декодировать JSON из тела запроса (email, password)
// --TODO: Проверить, что email и password не пустые
// --TODO: Вызвать h.AuthService.Login(email, password) — получить JWT токен
// --TODO: Если ошибка — вернуть 401 с сообщением "неправильный email или пароль"
// --TODO: Если успех — вернуть 200 с JSON: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
	username := UserInput.Username
	password := UserInput.Password

	if len(username) == 0 || len(password) == 0 {
		newErr := errors.New("заполните все поля")
		log.Printf("error while fill fields: %s", newErr)
		return
	}

	JWT, err := h.AuthService.Login(username, password)
	if err != nil {
		log.Printf("error while authorizate: %s", err)
		http.Error(w, "401 : Неверный логин или пароль", http.StatusUnauthorized)
		return
	}
	//установка куки
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    JWT, //сюда записать токен
		Path:     "/",
		HttpOnly: true,                    // Доступ только через HTTP, защита от XSS
		Secure:   true,                    // Только HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
		MaxAge:   900,                     // время жизни куки в секундах (поставил 15 минут)
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Удаление куки
	}
	http.SetCookie(w, cookie)
	h.AuthService.Logout()
	w.Write([]byte("Cookie deleted!"))
}

// GetBoards — обрабатывает GET /api/boards
// --TODO: Проверить, что метод GET (иначе 405)
// --TODO: Получить заголовок Authorization из r.Header
// --TODO: Проверить, что он начинается с "Bearer "
// --TODO: Извлечь токен — всё, что после "Bearer "
// --TODO: Вызвать h.AuthService.GetUserFromToken(token) — получить пользователя
// --TODO: Если токен невалиден — вернуть 401
// --TODO: Если токен валиден — получить доски через h.BoardService.GetBoards()
// TODO: Вернуть 200 с JSON: { "user": { "id", "email", "username" }, "boards": [...] }
func (h *Handler) GetBoards(w http.ResponseWriter, r *http.Request) {
	log.Println("GetBoard")
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		newBoard := new(models.Board)
		err := decoder.Decode(newBoard)
		if err != nil {
			log.Printf("error while unmarshalling JSON: %s", err)
			w.Write([]byte("{}"))
			return
		}
		err = h.BoardService.AddBoard(*newBoard)
		if err != nil {
			log.Printf("error while create Board: %s", err)
			http.Error(w, "400 : Bad Request", http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Cookie not found", http.StatusNotFound)
			return
		}

		tokenString := cookie.Value
		User, err := h.AuthService.GetUserFromToken(tokenString)
		//User, err := h.AuthService.GetCurrentUser()
		_ = User
		if err != nil {
			http.Error(w, "401 : Unauthorized or JWT not valid", http.StatusUnauthorized)
			log.Println("error:", err)
			return
		}
		Boards := h.BoardService.GetBoards()
		json_Boards, errB := json.Marshal(Boards)
		if errB != nil {
			log.Printf("error while serialize Boars: %s", errB)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_Boards)
	default:
		http.Error(w, "405 : NotAcceptable", http.StatusNotAcceptable)
		log.Printf("Запрос " + r.Method + ",а должен быть GET или POST")
		return
	}

}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}
	tokenString := cookie.Value
	user, err := h.AuthService.GetUserFromToken(tokenString)
	if err != nil {
		http.Error(w, "401 : Unauthorized or JWT not valid", http.StatusUnauthorized)
		log.Println("error:", err)
		return
	}

	user.Password = ""
	json_User, err := json.Marshal(user)
	if err != nil {
		log.Printf("error while marshalling User: %s", err)
		http.Error(w, "400 : Bad Request", http.StatusBadRequest)
		return
	}
	w.Write([]byte(json_User))
	log.Printf(string(json_User))
}

func (h *Handler) BoardDelete(w http.ResponseWriter, r *http.Request) {
	vars := r.PathValue("id")
	h.BoardService.DeleteBoard(vars)
}

func (h *Handler) BoardRestore(w http.ResponseWriter, r *http.Request) {
	vars := r.PathValue("boardId")
	h.BoardService.RestoreBoard(vars)
}
