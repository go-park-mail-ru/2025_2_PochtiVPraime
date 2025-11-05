package models

import (
	"time"
)

// Board — модель доски
type Board struct {
	ID        int64     `json:"id"`
	OwnerId   int64     `json:"owner_user_id"` // TODO: Добавить поле OwnerID — чтобы знать, кому принадлежит доска
	Title     string    `json:"title"`
	Image     string    `json:"image_id"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"` // TODO: Добавить CreatedAt
	UpdatedAt time.Time `json:"updated_at"` // TODO: Добавить UpdateAt
	// TODO: возможно ещё какие то поля
}

type BoardsData struct {
	ActiveBoards   []Board `json:"active_boards"`
	ArchivedBoards []Board `json:"archived_boards"`
}

type FullBoardData struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	//Image     string    `json:"image_id"`
	Lists []ListData `json:"lists"`
}
