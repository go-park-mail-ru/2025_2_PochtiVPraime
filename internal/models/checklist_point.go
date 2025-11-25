package models

import (
	"time"
)

type ChecklistPoint struct {
	ID          int64     `json:"id"`
	ChecklistID int64     `json:"checklist_id"`
	Content     string    `json:"content"`
	Checked     bool      `json:"checked"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
