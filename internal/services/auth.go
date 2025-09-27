package services

import (
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
)

// JWT_SECRET секрет для подписи токенов (в продакшене — из переменных окружения)
// TODO: Вынести в .env или переменные окружения (os.Getenv)
const JWT_SECRET = "super-secret-key-1234567890" // временно! Заменить на случайную строку!

// AuthService сервис для аутентификации и авторизации
// TODO: В будущем добавить:
// - db *sql.DB для хранения пользователей
// - logger *log.Logger для логирования событий
// - hasher *bcrypt.Hasher для хеширования паролей
type AuthService struct {
	// Поля будут добавлены позже пока пусто
}

// NewAuthService — конструктор для Dependency Injection
// TODO: В будущем принимать db, logger, hasher
// Сейчас — просто возвращаем пустой сервис
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register — регистрирует нового пользователя
// TODO: Проверить, что email не пустой
// TODO: Проверить, что email содержит "@"
// TODO: Проверить, что username не пустой
// TODO: Проверить, что password не короче 6 символов
// TODO: Проверить, что пользователь с таким email уже не существует
// TODO: Хешировать пароль
// TODO: Сохранить пользователя в базу данных (пока что в памяти)
// TODO: Вернуть *models.User без пароля
func (as *AuthService) Register(email, username, password string) (*models.User, error) {
	// Пока просто возвращаем nil — заглушка
	return nil, nil
}

// Login — авторизует пользователя и возвращает JWT токен
// TODO: Проверить, что email и password не пустые
// TODO: Найти пользователя по email
// TODO: Сравнить пароль (когда будем хешировать — использовать bcrypt.CompareHashAndPassword)
// TODO: Создать JWT токен с payload: { "userId": 123, "exp": 1720000000 }
// TODO: Вернуть токен и nil — если всё ок
// TODO: Вернуть ошибку "неправильный email или пароль" — если не найден
func (as *AuthService) Login(email, password string) (string, error) {
	// Пока просто возвращаем пустую строку — заглушка
	return "", nil
}

// GetUserFromToken — расшифровывает JWT и возвращает пользователя по ID
// TODO: Проверить, что токен не пустой
// TODO: Разобрать токен и проверить подпись (используя JWT_SECRET)
// TODO: Извлечь userID из payload
// TODO: Найти пользователя по userID
// TODO: Вернуть *User и nil — если токен валиден
// TODO: Вернуть nil и ошибку — если токен невалиден (истёк, подделан, отсутствует)
func (as *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// Пока просто возвращаем nil — заглушка
	return nil, nil
}
