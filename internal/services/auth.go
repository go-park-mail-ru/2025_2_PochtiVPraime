package services

import (
	"errors"
	"log"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"golang.org/x/crypto/bcrypt"
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

var userId int = 0
var storeUsers map[string]models.User = map[string]models.User{}

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
	if len(storeBoards) == 0 {
		userId = 0
	} else {
		userId++
	}
	cost := bcrypt.DefaultCost
	encode_pass, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Printf("error while encode password: %s", err)
		return nil, err
	}
	storeUsers[username] = models.User{ID: userId, Email: email, Username: username, Password: string(encode_pass)}
	newUser := storeUsers[username]
	return &newUser, nil
}

// Login — авторизует пользователя и возвращает JWT токен
// TODO: Проверить, что email и password не пустые
// --TODO: Найти пользователя по email
// --TODO: Сравнить пароль (когда будем хешировать — использовать bcrypt.CompareHashAndPassword)
// TODO: Создать JWT токен с payload: { "userId": 123, "exp": 1720000000 }
// TODO: Вернуть токен и nil — если всё ок
// TODO: Вернуть ошибку "неправильный email или пароль" — если не найден
func (as *AuthService) Login(username, password string) (string, error) {
	// Пока просто возвращаем пустую строку — заглушка
	User, flag := storeUsers[username]
	if !flag {
		log.Printf("wrong username")
		return "", errors.New("Нет пользователя с таким username")
	}
	err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(password))
	if err != nil {
		log.Printf("Wrong password: %s", err)
		return "", err
	}

	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil
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
