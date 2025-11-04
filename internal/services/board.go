package services

import (
	"context"
	"errors"

	repository "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/Repository"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
)

// BoardService — сервис для работы с досками
type BoardService struct {
	BoardRepository repository.BoardsRepository
	UserRepository  repository.UserRepository
	// Здесь будут зависимости в будущем
}

// NewBoardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewBoardService(boardRepository *repository.BoardsRepository) *BoardService {
	return &BoardService{
		BoardRepository: *boardRepository,
	}
}

// GetBoards — возвращает список досок
// TODO: Получить доски только для авторизованного пользователя (по userId)
// TODO: Загружать доски из базы данных (ну или пока что просто из списка)
func (bs *BoardService) GetBoards(ctx context.Context, userId int64) (*models.BoardsData, error) {
	var userBoards = models.BoardsData{}
	var rawUserBoards, err = bs.BoardRepository.GetBoardsByOwner(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, value := range rawUserBoards {
		if !value.Archived {
			userBoards.ActiveBoards = append(userBoards.ActiveBoards, *value)
		} else {
			userBoards.ArchivedBoards = append(userBoards.ArchivedBoards, *value)
		}
	}
	return &userBoards, nil
}

// AddBoard — создаёт новую доску, не обязательно, но как будто бы надо
// --TODO: Проверить, что title не пустой
// TODO: Проверить, что пользователь авторизован (через session)
// TODO: Сохранить доску в БД с привязкой к userId
// TODO: Вернуть созданную доску
func (bs *BoardService) AddBoard(ctx context.Context, board *models.Board) error {

	if board.Title == "" {
		return errors.New("Нет Title")
	}
	board.Archived = false
	if currentUser.Email == "" {
		return errors.New("Пользователь не авторизирован")
	}
	_, err := bs.BoardRepository.CreateBoard(ctx, board)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BoardService) DeleteBoard(ctx context.Context, boardId int64) error {
	// Пока просто возвращаем nil
	err := bs.BoardRepository.DeleteBoard(ctx, boardId)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BoardService) RestoreBoard(ctx context.Context, boardId int64) error {
	err := bs.BoardRepository.RestoreBoard(ctx, boardId)
	if err != nil {
		return err
	}
	return nil
}
