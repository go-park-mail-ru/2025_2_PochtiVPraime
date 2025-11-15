package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type BoardMemberRepository interface {
	CreateBoardMember(ctx context.Context, user *models.BoardMember) (*models.BoardMember, error)
	GetBoardMemberById(ctx context.Context, id int64) (*models.BoardMember, error)
	GetBoardMemberByUserId(ctx context.Context, userId int64) (*models.BoardMember, error)
	ChangeRole(ctx context.Context, newRole string, memberId int64) error
	DeleteBoardMember(ctx context.Context, id int64) error
}

type BoardsRepository interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetBoardById(ctx context.Context, id int64) (*models.Board, error)
	GetBoardsByOwner(ctx context.Context, ownerID int64) ([]*models.Board, error)
	UpdateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	ArchiveBoard(ctx context.Context, id int64) error
	RestoreBoard(ctx context.Context, id int64) error
	DeleteBoard(ctx context.Context, id int64) error
}

type ListsRepository interface {
	SaveList(ctx context.Context, list *models.List) (*models.List, error)
	GetListByID(ctx context.Context, id int64) (*models.List, error)
	GetListsByBoardID(ctx context.Context, id int64) ([]*models.List, error)
	UpdateList(ctx context.Context, list *models.List) (*models.List, error)
	DeleteList(ctx context.Context, id int64) error
}

type CardsRepository interface {
	CreateCard(ctx context.Context, card *models.Card) (*models.Card, error)
	GetCard(ctx context.Context, id int64) (*models.Card, error)
	GetCardsByList(ctx context.Context, listID int64) ([]*models.Card, error)
	UpdateCard(ctx context.Context, card *models.Card) (*models.Card, error)
	DeleteCard(ctx context.Context, id int64) error
	UpdateCardPosition(ctx context.Context, cardID int64, newPosition int, newListID int64) error
	GetCardsByBoardMember(ctx context.Context, boardMemberID int64) ([]*models.Card, error)
}

type SupportRepository interface {
	CreateSupportForm(ctx context.Context, board *models.SupportForm) error
	GetSupportFormById(ctx context.Context, id int64) (*models.SupportForm, error)
	GetSupportFormsByOwner(ctx context.Context, ownerID int64) ([]*models.SupportForm, error)
	GetAllSupportForms(ctx context.Context) ([]*models.SupportForm, error)
	//ArchiveBoard(ctx context.Context, id int64) error
	//RestoreBoard(ctx context.Context, id int64) error
	DeleteSupportForm(ctx context.Context, id int64) error
}
