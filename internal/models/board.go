package models

// Board — модель доски
type Board struct {
	Id        string `json:"id"`
	OwnerId   int    `json:"ownerId"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Archived  bool   `json:"archived"`
	CreatedAt string `json:"createAt"`
	// TODO: Добавить поле UserID — чтобы знать, кому принадлежит доска
	// TODO: Добавить CreatedAt
	// TODO: возможно ещё какие то поля
}

type BoardsData struct {
	ActiveBoards   []Board `json:"activeBoards"`
	ArchivedBoards []Board `json:"archivedBoards"`
}
