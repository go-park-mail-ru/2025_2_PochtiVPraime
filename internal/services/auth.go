package services

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
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
	UserRepository repository.UserRepository
	// Поля будут добавлены позже пока пусто
}

// NewAuthService — конструктор для Dependency Injection
// TODO: В будущем принимать db, logger, hasher
// Сейчас — просто возвращаем пустой сервис
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepo,
	}
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
func (as *AuthService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	email := user.Email
	if !strings.Contains(email, "@") || len(email) == 0 { //наверное len должна быть хотябы 6
		newErr := errors.New("Не содержит @ или слишком короткий email") //тк len(a@b.ru)
		log.Printf("error while email not valid: %s", newErr)
		return nil, newErr
	}
	username := user.Username
	if len(username) <= 0 || len(username) > 25 {
		newErr := errors.New("слишком короткое или слишком длинное имя")
		log.Printf("error while name not valid: %s", newErr)
		return nil, newErr
	}

	password := user.Password
	if len(password) < 6 {
		newErr := errors.New("слишком короткий пароль")
		log.Printf("error while name not valid: %s", newErr)
		return nil, newErr
	}

	cost := bcrypt.DefaultCost
	encode_pass, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	user.Password = string(encode_pass)
	if err != nil {
		log.Printf("error while encode password: %s", err)
		return nil, err
	}
	/*
		as.UserRepository.
		if err != nil {
			log.Printf(" Такое имя уже занято")
			return nil, errors.New("Такое имя уже занято")
		}

		for _, user := range storeUsers {
			if email == user.Email {
				log.Printf("Пользователь с таким email уже существует")
				return nil, errors.New("Пользователь с таким email уже существует")
			}
		}
	*/
	//storeUsers[username] = models.User{ID: userId, Email: email, Username: username, Password: string(encode_pass)}
	//newUser := storeUsers[username]
	user, err = as.UserRepository.CreateUser(ctx, user)
	if err != nil {
		log.Printf("error while saving User in DB: %s", err)
		return nil, err
	}
	//validUser.Password = ""
	//log.Println(storeUsers[username])
	//userId++
	return user, nil
}

// Login — авторизует пользователя и возвращает JWT токен
// --TODO: Проверить, что email и password не пустые
// --TODO: Найти пользователя по email
// --TODO: Сравнить пароль (когда будем хешировать — использовать bcrypt.CompareHashAndPassword)
// --TODO: Создать JWT токен с payload: { "userId": 123, "exp": 1720000000 }
// --TODO: Вернуть токен и nil — если всё ок
// --TODO: Вернуть ошибку "неправильный email или пароль" — если не найден
func (as *AuthService) Login(ctx context.Context, user *models.User) (string, error) {
	username := user.Username
	password := user.Password

	if len(username) == 0 || len(password) == 0 {
		newErr := errors.New("заполните все поля")
		log.Printf("error while fill fields: %s", newErr)
		return "", newErr
	}
	log.Println(username + " " + password)
	localUser, err := as.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("wrong username")
		return "", errors.New("Нет пользователя с таким именем")
	}
	err = bcrypt.CompareHashAndPassword([]byte(localUser.Password), []byte(password))
	if err != nil {
		log.Printf("Wrong password: %s", err)
		return "", errors.New("Неправильный пароль")
	}
	claims := jwt.MapClaims{
		"userId": localUser.ID,
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
func (as *AuthService) GetUserFromToken(ctx context.Context, tokenString string) (*models.User, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Парсим токен без проверки подписи (только для получения claims)
	//token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method for JWT token")
			return nil, errors.New("unexpected signing method for JWT token")

		}

		return []byte(JWT_SECRET), nil

	})

	if err != nil {

		log.Println("JWT parsing error:", err)

		return nil, errors.New("JWT parsing error:")

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("invalid token")
		return nil, errors.New("invalid token")
	}

	// Получаем user_id из claims
	id, ok := claims["userId"]
	if !ok {
		log.Println("user_id не найден в токен")
		return nil, errors.New("user_id не найден в токене")
	}

	var user *models.User
	currentId, ok := id.(float64)
	log.Println(id)
	if !ok {
		log.Println("не смог привести user_id к int")
		return nil, errors.New("не смог привести user_id к int")
	}
	newId := int64(currentId)
	log.Println(newId)
	user, err = as.UserRepository.GetUserByID(ctx, newId)
	log.Println(user)
	return user, nil
}

func (as *AuthService) Logout() {
	//currentUser = models.User{}
}

func (as *AuthService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	localUser, err := as.UserRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		log.Printf("wrong username")
		return nil, errors.New("Нет пользователя с таким именем")
	}
	if user.Email != "" {
		localUser.Email = user.Email
	}
	if user.Username != "" {
		localUser.Username = user.Username
	}
	localUser.UpdatedAt = time.Now()
	//изменение авы в будущем
	return as.UserRepository.UpdateUser(ctx, localUser)
}

func (as *AuthService) PasswordUpdate(ctx context.Context, oldPassword string, newPassword string, userId int64) (*models.User, error) {
	localUser, err := as.UserRepository.GetUserByID(ctx, userId)
	if err != nil {
		log.Printf("wrong username")
		return nil, errors.New("Нет пользователя с таким именем")
	}
	localUser.UpdatedAt = time.Now()
	if oldPassword == newPassword {
		return nil, errors.New("Новый и старый пароли не должны совпадать")
	}
	return as.UserRepository.UpdateUser(ctx, localUser)
}
