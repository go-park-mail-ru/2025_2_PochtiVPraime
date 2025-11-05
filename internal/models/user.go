package models

import "time"

// User — модель пользователя (данные, которые мы храним и передаём)
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"` // TODO: Добавить поле Password (но не отправлять в ответе через json:"-")
	AvatarID  int64     `json:"avatar_id,omitempty"`
	CreatedAt time.Time `json:"created_at"` //? TODO: Добавить CreatedAt для регистрации
	UpdatedAt time.Time `json:"updated_at"` //? TODO: Добавить UpdatedAt для изменений
}
