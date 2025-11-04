package models

import "time"

// Board — модель доски
type Board struct {
	ID        int64     `json:"id"`
	OwnerId   int64     `json:"owner_id"` // TODO: Добавить поле OwnerID — чтобы знать, кому принадлежит доска
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"create_at"` // TODO: Добавить CreatedAt
	// TODO: возможно ещё какие то поля
}

type BoardsData struct {
	ActiveBoards   []Board `json:"active_boards"`
	ArchivedBoards []Board `json:"archived_boards"`
}
