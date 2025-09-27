package models

// User — модель пользователя (данные, которые мы храним и передаём)
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	// TODO: Добавить поле Password (но не отправлять в ответе через json:"-")
	// TODO: Добавить CreatedAt для регистрации
	// TODO: Добавить UpdatedAt для изменений
}
