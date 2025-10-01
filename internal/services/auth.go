package services

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT_SECRET секрет для подписи токенов (в продакшене — из переменных окружения)
// TODO: Вынести в .env или переменные окружения (os.Getenv)
var JWT_SECRET = []byte("super-secret-key-1234567890") // временно! Заменить на случайную строку!

// AuthService сервис для аутентификации и авторизации
// TODO: В будущем добавить:
// - db *sql.DB для хранения пользователей
// - logger *log.Logger для логирования событий
// - hasher *bcrypt.Hasher для хеширования паролей
type AuthService struct {
	// Поля будут добавлены позже пока пусто
}

var currentUser models.User

var userId int = 0
var storeUsers = map[string]models.User{}

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
// --TODO: Хешировать пароль
// --TODO: Сохранить пользователя в базу данных (пока что в памяти)
// --TODO: Вернуть *models.User без пароля
func (as *AuthService) Register(email, username, password string) (*models.User, error) {
	if len(storeUsers) == 0 {
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
	_, flag := storeUsers[username]
	if flag {
		log.Printf(" user already exist: %s", err)
		return nil, err
	}

	storeUsers[username] = models.User{ID: userId, Email: email, Username: username, Password: string(encode_pass)}
	newUser := storeUsers[username]
	log.Println(storeUsers[username])
	//log.Println(len(storeUsers))
	return &newUser, nil
}

// Login — авторизует пользователя и возвращает JWT токен
// --TODO: Проверить, что email и password не пустые
// --TODO: Найти пользователя по email
// --TODO: Сравнить пароль (когда будем хешировать — использовать bcrypt.CompareHashAndPassword)
// --TODO: Создать JWT токен с payload: { "userId": 123, "exp": 1720000000 }
// --TODO: Вернуть токен и nil — если всё ок
// --TODO: Вернуть ошибку "неправильный email или пароль" — если не найден
func (as *AuthService) Login(username, password string) (string, error) {
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
	currentUser = User
	claims := jwt.MapClaims{
		"userId": storeUsers[username].ID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JWT_SECRET)
}

// GetUserFromToken — расшифровывает JWT и возвращает пользователя по ID
// TODO: Проверить, что токен не пустой
// TODO: Разобрать токен и проверить подпись (используя JWT_SECRET)
// TODO: Извлечь userID из payload
// TODO: Найти пользователя по userID
// TODO: Вернуть *User и nil — если токен валиден
// TODO: Вернуть nil и ошибку — если токен невалиден (истёк, подделан, отсутствует)
func (as *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Парсим токен без проверки подписи (только для получения claims)
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, errors.New("jwt.NewParser().ParseUnverified")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("неверный формат claims")
	}

	// Получаем user_id из claims
	id, ok := claims["userId"]
	if !ok {
		return nil, errors.New("user_id не найден в токене")
	}

	var user models.User
	currentId, ok := id.(float64)
	if !ok {
		return nil, errors.New("не смог привести user_id к int")
	}

	for key, value := range storeUsers {
		if float64(value.ID) == currentId {
			user = storeUsers[key]
		}
	}
	log.Println(user)
	return &user, nil
}

func (as *AuthService) GetCurrentUser() (*models.User, error) {
	if currentUser.Email == "" {
		return nil, errors.New("Не авторизирован")
	}
	return &currentUser, nil
}

func (as *AuthService) LogoutUser() {
	currentUser = models.User{}
}
