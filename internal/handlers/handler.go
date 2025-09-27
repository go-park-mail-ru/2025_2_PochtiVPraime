package handlers

import (
	"net/http"

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

// Register — обрабатывает POST /register
// TODO: Проверить, что метод POST (иначе 405)
// TODO: Декодировать JSON из тела запроса (email, username, password)
// TODO: Проверить, что email не пустой и содержит "@"
// TODO: Проверить, что username не пустой и не слишком длинный
// TODO: Проверить, что password не короче 6 символов
// TODO: Проверить, что пользователь с таким email уже не существует
// TODO: Вызвать h.AuthService.Register(email, username, password)
// TODO: Если ошибка — вернуть 400 с сообщением об ошибке
// TODO: Если успех — вернуть 201 с JSON: { "user": { "id", "email", "username" } }
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Пока просто отвечаем заглушкой
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "register работает"}`))
}

// Login — обрабатывает POST /login
// TODO: Проверить, что метод POST (иначе 405)
// TODO: Декодировать JSON из тела запроса (email, password)
// TODO: Проверить, что email и password не пустые
// TODO: Вызвать h.AuthService.Login(email, password) — получить JWT токен
// TODO: Если ошибка — вернуть 401 с сообщением "неправильный email или пароль"
// TODO: Если успех — вернуть 200 с JSON: { "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." }
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Пока просто отвечаем заглушкой
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "login работает"}`))
}

// GetBoards — обрабатывает GET /get-boards
// TODO: Проверить, что метод GET (иначе 405)
// TODO: Получить заголовок Authorization из r.Header
// TODO: Проверить, что он начинается с "Bearer "
// TODO: Извлечь токен — всё, что после "Bearer "
// TODO: Вызвать h.AuthService.GetUserFromToken(token) — получить пользователя
// TODO: Если токен невалиден — вернуть 401
// TODO: Если токен валиден — получить доски через h.BoardService.GetBoards()
// TODO: Вернуть 200 с JSON: { "user": { "id", "email", "username" }, "boards": [...] }
func (h *Handler) GetBoards(w http.ResponseWriter, r *http.Request) {
	// Пока просто отвечаем заглушкой
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"user": null, "boards": []}`))
}
