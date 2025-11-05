package models

import "time"

//модель карточки
type Card struct {
	ID                  int64     `json:"id"`
	AuthorBoardMemberId int64     `json:"author_board_member_id"`
	ListId              int64     `json:"list_id"`
	Content             string    `json:"content"`
	Position            int       `json:"position"`
	Completed           bool      `json:"completed"`
	CreatedAt           time.Time `json:"created_at"` //?
	UpdatedAt           time.Time `json:"updated_at"` //?
	CompleteBefore      time.Time `json:"complete_before"`
}

type CardData struct {
	ID        int64  `json:"id"`
	Content   string `json:"text"`
	Completed bool   `json:"isCompleted"`
}
