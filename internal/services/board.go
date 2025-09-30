package services

import "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"

var boardsId int
var storeBoards map[int]models.Board

var test = models.BoardsData{
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
			OwnerId:   1,
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
		OwnerId:   1,
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
			OwnerId:   1,
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
	// Пока возвращаем пустой список
	return test
}

// AddBoard — создаёт новую доску, не обязательно, но как будто бы надо
// TODO: Проверить, что name не пустой
// TODO: Проверить, что пользователь авторизован (через session)
// TODO: Сохранить доску в БД с привязкой к userId
// TODO: Вернуть созданную доску
func (bs *BoardService) AddBoard(name string) (*models.Board, error) {
	// Пока просто возвращаем nil

	return nil, nil
}
