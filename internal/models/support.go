package models

import "time"

// User — модель пользователя (данные, которые мы храним и передаём)
type SupportForm struct {
	ID           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	HelperId     int64     `json:"helper_id"`
	FormType     string    `json:"form_type"`
	FormStatus   string    `json:"form_status"` // TODO: Добавить поле Password (но не отправлять в ответе через json:"-")
	Text         string    `json:"text"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"` //? TODO: Добавить CreatedAt для регистрации
	UpdatedAt    time.Time `json:"updated_at"` //? TODO: Добавить UpdatedAt для изменений
}
