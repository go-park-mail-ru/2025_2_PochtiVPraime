package models

import "time"

//модель карточки
type Card struct {
	ID                  int64     `json:"id,omitempty"`
	AuthorBoardMemberId int64     `json:"author_board_member_id,omitempty"`
	ListId              int64     `json:"list_id,omitempty"`
	Content             string    `json:"content,omitempty"`
	Position            int       `json:"position,omitempty"`
	Completed           bool      `json:"completed,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"` //?
	UpdatedAt           time.Time `json:"updated_at,omitempty"` //?
	CompleteBefore      time.Time `json:"complete_before,omitempty"`
}

type CardData struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"isCompleted"`
}
