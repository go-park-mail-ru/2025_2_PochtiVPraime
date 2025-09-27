package models

// Board — модель доски
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// TODO: Добавить поле UserID — чтобы знать, кому принадлежит доска
	// TODO: Добавить CreatedAt
	// TODO: возможно ещё какие то поля
}
