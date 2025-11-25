package middleware

import (
	"context"
	"net/http"
	"strings"
)

// ключ для хранения UserId в контексте
type contextKey string

const (
	TokenKey contextKey = "tokenString"
)

// AuthMiddleware извлекает UserId из cookie и добавляет в контекст
func FillContext(r *http.Request) context.Context {
	// Извлекаем куки
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return r.Context()
	}

	// Очищаем и валидируем значение
	tokenString := strings.TrimSpace(cookie.Value)

	// Добавляем userID в контекст
	return context.WithValue(r.Context(), TokenKey, tokenString)
}
