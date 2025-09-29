package repository

import "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"

type UserRepository interface {
	SaveUser(user models.User) error
	Authorizate(email string, password string) (user models.User)
	FindByID(id int) (models.User, error)
	GetUserBoard(userId int)
	UpdateUser(newUser models.User) (models.User, error)
}

type BoardsRepository interface {
	AddBoard(board models.Board) error
	GetBoardById(id int) ([](models.Board), error)
	UpdateBoard(newBoard models.Board) (models.Board, error)
}
