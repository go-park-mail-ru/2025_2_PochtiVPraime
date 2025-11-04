package models

import "time"

type attachment struct {
	ID        int64     `json:"id"`
	cardId    int64     `json:"card_id"`
	Title     string    `json:"title"`
	FileUrl   string    `json:"file_url"`
	Position  int       `json:"position"`
	createdAt time.Time `json:"created_at"` //?
	updatedAt time.Time `json:"updated_at"` //?
}
