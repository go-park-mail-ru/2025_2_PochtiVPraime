package models

import (
	"time"
)

type List struct {
	ID        int64     `json:"id"`
	BoardId   int64     `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListData struct {
	ID    int64      `json:"id"`
	Title string     `json:"title"`
	Tasks []CardData `json:"tasks"`
}
