package models

import (
	"time"
)

// Board — модель доски
type Checklist struct {
	ID        int64     `json:"id"`
	CardId    int64     `json:"card_id"` // TODO: Добавить поле OwnerID — чтобы знать, кому принадлежит доска
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"create_at"` // TODO: Добавить CreatedAt
	UpdatedAt time.Time `json:"updated_at"`
}

type ChecklistData struct {
	CardId          int64            `json:"card_id"` // TODO: Добавить поле OwnerID — чтобы знать, кому принадлежит доска
	Title           string           `json:"title"`
	ChecklistPoints []ChecklistPoint `json:"checklist_points"`
	CreatedAt       time.Time        `json:"create_at"` // TODO: Добавить CreatedAt
	UpdatedAt       time.Time        `json:"updated_at"`
}
