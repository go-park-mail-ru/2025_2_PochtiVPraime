package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	repository "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/Repository"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
)

// ListService — сервис для работы со списками карточек
type ListService struct {
	ListRepository  repository.ListsRepository
	BoardRepository repository.BoardsRepository
	CardsRepository repository.CardsRepository
	// Здесь будут зависимости в будущем
}

// NewListService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewListService(listRepository *repository.ListsRepository, boardRepository *repository.BoardsRepository, cardRepository *repository.CardsRepository) *ListService {
	return &ListService{
		ListRepository:  *listRepository,
		BoardRepository: *boardRepository,
		CardsRepository: *cardRepository,
	}
}

// CreteList - создаёт список карточек
func (ls *ListService) CreateList(ctx context.Context, title string, boardId int64, userId int64) (*models.List, error) {
	// Проверяем существование доски
	board, err := ls.BoardRepository.GetBoardById(ctx, boardId)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	// Проверяем, что пользователь является владельцем доски
	if board.OwnerId != userId {
		return nil, fmt.Errorf("user not owner of board: %w", err)
	}

	// Получаем текущие списки доски для определения позиции
	existingLists, err := ls.ListRepository.GetListsByBoardID(ctx, boardId)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing lists: %w", err)
	}

	// Определяем позицию нового списка (в конец)
	position := 0
	if len(existingLists) > 0 {
		position = existingLists[len(existingLists)-1].Position + 1
	}

	// Создаем объект списка
	list := &models.List{
		BoardId:   boardId,
		Title:     title,
		Position:  position,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем в репозитории
	createdList, err := ls.ListRepository.SaveList(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return createdList, nil
}

// GetLists - получаем списки карточек по id доски
func (ls *ListService) GetLists(ctx context.Context, boardId int64) ([]*models.List, error) {
	/*
		// Проверяем права доступа к доске
		board, err := ls.BoardRepository.GetBoardById(ctx, boardId)
		if err != nil {
			return nil, fmt.Errorf("failed to get board: %w", err)
		}

			if board.OwnerID != userID {
				return nil, ErrListAccessDenied
			}
	*/

	// Получаем списки доски
	lists, err := ls.ListRepository.GetListsByBoardID(ctx, boardId)
	if err != nil {
		return nil, fmt.Errorf("failed to get board lists: %w", err)
	}

	return lists, nil
}

// UpdateList - изменяет карточку и возвращает обновленную версию
func (ls *ListService) UpdateList(ctx context.Context, newList *models.List, userId int64) (*models.List, error) {
	// Получаем текущий список
	list, err := ls.ListRepository.GetListByID(ctx, newList.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	// Получаем доску
	board, err := ls.BoardRepository.GetBoardById(ctx, newList.BoardId)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	if board.OwnerId != userId {
		return nil, errors.New("access denied to board")
	}

	list.UpdatedAt = time.Now()

	// Сохраняем изменения
	updatedList, err := ls.ListRepository.UpdateList(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return updatedList, nil
}

// DeleteList - удаляет список карточек по id списка
func (ls *ListService) DeleteListWithCard(ctx context.Context, listId int64, userId int64) error {
	// Получаем список для проверки прав
	list, err := ls.ListRepository.GetListByID(ctx, listId)
	if err != nil {
		return fmt.Errorf("failed to get list: %w", err)
	}

	// Проверяем права доступа к доске списка
	board, err := ls.BoardRepository.GetBoardById(ctx, list.BoardId)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}

	if board.OwnerId != userId {
		return errors.New("access denied to board")
	}

	// Удаляем список (в репозитории должна быть каскадная обработка)
	err = ls.ListRepository.DeleteList(ctx, listId)
	if err != nil {
		return fmt.Errorf("failed to delete list with cards: %w", err)
	}

	return nil
}
