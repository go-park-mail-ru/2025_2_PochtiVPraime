package services

import (
	"context"
	"errors"
	"log"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

// BoardService — сервис для работы с досками
type BoardService struct {
	ListRepository  repository.ListsRepository
	CardRepository  repository.CardsRepository
	BoardRepository repository.BoardsRepository
	UserRepository  repository.UserRepository
	// Здесь будут зависимости в будущем
}

// NewBoardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewBoardService(boardRepository repository.BoardsRepository, listRepository repository.ListsRepository,
	cardRepository repository.CardsRepository, userRepository repository.UserRepository) *BoardService {
	return &BoardService{
		BoardRepository: boardRepository,

		ListRepository: listRepository,
		CardRepository: cardRepository,
		UserRepository: userRepository,
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

func (bs *BoardService) GetBoard(ctx context.Context, boardId int64) (*models.FullBoardData, error) {
	// 1. Получаем базовую информацию о доске
	board, err := bs.BoardRepository.GetBoardById(ctx, boardId)
	if err != nil {
		log.Printf("Error getting board by id %d: %v", boardId, err)
		return nil, err
	}
	// Получаем все списки для этой доски
	lists, err := bs.ListRepository.GetListsByBoardID(ctx, boardId)
	if err != nil {
		log.Printf("Error getting lists for board %d: %v", boardId, err)
		return nil, err
	}

	// 3. Для каждого списка получаем карточки
	listData := make([]models.ListData, 0, len(lists))
	for _, list := range lists {
		cards, err := bs.CardRepository.GetCardsByList(ctx, list.ID)
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
	// 6. Собираем полный объект доски
	fullBoard := &models.FullBoardData{
		ID:    board.ID,
		Title: board.Title,
		Lists: listData,
	}

	return fullBoard, nil
}
