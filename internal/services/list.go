package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

// ListService — сервис для работы со списками карточек
type ListService struct {
	ListRepository  repository.ListsRepository
	BoardRepository repository.BoardsRepository
	CardsRepository repository.CardsRepository
	// Здесь будут зависимости в будущем
}

// NewListService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewListService(listRepository repository.ListsRepository, boardRepository repository.BoardsRepository, cardRepository repository.CardsRepository) *ListService {
	return &ListService{
		ListRepository:  listRepository,
		BoardRepository: boardRepository,
		CardsRepository: cardRepository,
	}
}

// CreteList - создаёт список карточек
func (ls *ListService) CreateList(ctx context.Context, title string, boardId int64, userId int64) (*models.List, error) {
	// Проверяем существование доски
	board, err := ls.BoardRepository.GetBoardById(ctx, boardId)
	_ = board
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	/*
		// Проверяем, что пользователь является владельцем доски
		if board.OwnerId != userId {
			return nil, fmt.Errorf("user not owner of board: %w", err)
		}
	*/
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
func (ls *ListService) GetLists(ctx context.Context, boardId int64) ([]models.ListData, error) {
	// Получаем все списки для этой доски
	lists, err := ls.ListRepository.GetListsByBoardID(ctx, boardId)
	if err != nil {
		log.Printf("Error getting lists for board %d: %v", boardId, err)
		return nil, err
	}

	// 3. Для каждого списка получаем карточки
	listData := make([]models.ListData, 0, len(lists))
	for _, list := range lists {
		cards, err := ls.CardsRepository.GetCardsByList(ctx, list.ID)
		if err != nil {
			log.Printf("Error getting cards for list %d: %v", list.ID, err)
			continue // Продолжаем обработку других списков
		}

		// 4. Преобразуем карточки в CardData
		cardData := make([]models.CardData, 0, len(cards))
		for _, card := range cards {
			cardData = append(cardData, models.CardData{
				ID:        card.ID,
				Content:   card.Content,
				Completed: card.Completed,
			})
		}

		// 5. Создаем ListData
		listData = append(listData, models.ListData{
			ID:    list.ID,
			Title: list.Title,
			Tasks: cardData,
		})
	}
	return listData, nil
}

// UpdateList - изменяет карточку и возвращает обновленную версию
func (ls *ListService) UpdateList(ctx context.Context, newList *models.List, userId int64) (*models.List, error) {
	// Получаем текущий список
	list, err := ls.ListRepository.GetListByID(ctx, newList.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	/*
		// Получаем доску
		board, err := ls.BoardRepository.GetBoardById(ctx, newList.BoardId)
		if err != nil {
			return nil, fmt.Errorf("failed to get board: %w", err)
		}

		if board.OwnerId != userId {
			return nil, errors.New("access denied to board")
		}
	*/
	list.UpdatedAt = time.Now()
	list.Title = newList.Title
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

func (ls *ListService) GetList(ctx context.Context, listId int64) (*models.ListData, error) {
	// Для списка получаем карточки
	list, err := ls.ListRepository.GetListByID(ctx, listId)
	if err != nil {
		log.Printf("Error getting list for listId %d: %v", listId, err)
		return nil, err
	}
	cards, err := ls.CardsRepository.GetCardsByList(ctx, list.ID)
	if err != nil {
		log.Printf("Error getting cards for list %d: %v", list.ID, err)
		return nil, err
	}

	// Преобразуем карточки в CardData
	cardData := make([]models.CardData, 0, len(cards))
	for _, card := range cards {
		cardData = append(cardData, models.CardData{
			ID:        card.ID,
			Content:   card.Content,
			Completed: card.Completed,
		})
	}
	listData := &models.ListData{
		ID:    list.ID,
		Title: list.Title,
		Tasks: cardData,
	}
	return listData, nil
}
