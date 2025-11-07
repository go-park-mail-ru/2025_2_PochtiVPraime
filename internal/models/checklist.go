package models

// Board — модель доски
type Checklist struct {
	ID        int64  `json:"id"`
	CardId    int64  `json:"card_id"` // TODO: Добавить поле OwnerID — чтобы знать, кому принадлежит доска
	Title     string `json:"title"`
	CreatedAt string `json:"create_at"` // TODO: Добавить CreatedAt
	UpdatedAt string `json:"updated_at"`
	// TODO: возможно ещё какие то поля
}
