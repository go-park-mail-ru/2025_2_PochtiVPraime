package services

import (
	"errors"
	"slices"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
)

var boardsId int
var storeBoards = models.BoardsData{
	ActiveBoards: []models.Board{
		{
			Id:        "board_1",
			OwnerId:   1,
			Title:     "Планы на год",
			Image:     "",
			Archived:  false,
			CreatedAt: " ",
		},
		{
			Id:        "board_2",
			OwnerId:   2,
			Title:     "Рабочие задачи",
			Image:     "/Images/default-board-bg.jpg",
			Archived:  false,
			CreatedAt: "",
		},
		{
			Id:        "board_3",
			OwnerId:   1,
			Title:     "Личные цели",
			Image:     "/Images/default-board-bg.jpg",
			Archived:  false,
			CreatedAt: "",
		},
	},

	ArchivedBoards: []models.Board{{
		Id:        "board_4",
		OwnerId:   2,
		Title:     "Проект А",
		Image:     "/Images/default-board-bg.jpg",
		Archived:  true,
		CreatedAt: "",
	},
		{
			Id:        "board_5",
			OwnerId:   1,
			Title:     "Проект Б",
			Image:     "/Images/default-board-bg.jpg",
			Archived:  true,
			CreatedAt: "",
		},
		{
			Id:        "board_6",
			OwnerId:   2,
			Title:     "Идеи",
			Image:     "/Images/default-board-bg.jpg",
			Archived:  true,
			CreatedAt: "",
		}},
}

// BoardService — сервис для работы с досками
type BoardService struct {
	// Здесь будут зависимости в будущем
}

// NewBoardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewBoardService() *BoardService {
	return &BoardService{}
}

// GetBoards — возвращает список досок
// TODO: Получить доски только для авторизованного пользователя (по userId)
// TODO: Загружать доски из базы данных (ну или пока что просто из списка)
func (bs *BoardService) GetBoards() models.BoardsData {
	var userBoards = models.BoardsData{}
	for _, value := range storeBoards.ActiveBoards {
		if value.OwnerId == currentUser.ID {
			userBoards.ActiveBoards = append(userBoards.ActiveBoards, value)
		}
	}
	for _, value := range storeBoards.ArchivedBoards {
		if value.OwnerId == currentUser.ID {
			userBoards.ArchivedBoards = append(userBoards.ArchivedBoards, value)
		}
	}
	return userBoards
}

// AddBoard — создаёт новую доску, не обязательно, но как будто бы надо
// --TODO: Проверить, что title не пустой
// TODO: Проверить, что пользователь авторизован (через session)
// TODO: Сохранить доску в БД с привязкой к userId
// TODO: Вернуть созданную доску
func (bs *BoardService) AddBoard(board models.Board) error {

	if board.Title == "" {
		return errors.New("Нет Title")
	}
	board.Archived = false
	if currentUser.Email == "" {
		return errors.New("Пользователь не авторизирован")
	}
	board.OwnerId = currentUser.ID
	storeBoards.ActiveBoards = append(storeBoards.ActiveBoards, board)
	return nil
}

func (bs *BoardService) DeleteBoard(boardId string) error {
	// Пока просто возвращаем nil
	for key, value := range storeBoards.ArchivedBoards {
		if value.Id == boardId {
			storeBoards.ArchivedBoards = slices.Delete(storeBoards.ArchivedBoards, key, key+1)
			return nil
		}
	}
	return errors.New("Такой доски не существует в ArchivedBoards")
}

func (bs *BoardService) RestoreBoard(boardId string) error {
	for key, value := range storeBoards.ArchivedBoards {
		if value.Id == boardId {
			value.Archived = false
			storeBoards.ActiveBoards = append(storeBoards.ActiveBoards, value)
			storeBoards.ArchivedBoards = slices.Delete(storeBoards.ArchivedBoards, key, key+1)
			return nil
		}
	}
	return errors.New("Не удалось восстановить доску")
}
