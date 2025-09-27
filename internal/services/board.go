package services

import "github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"

// BoardService — сервис для работы с досками
type BoardService struct {
	// Здесь будут зависимости в будущем
}

// NewBoardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewBoardService() *BoardService {
	return &BoardService{}
}

// GetBoards — возвращает список досок
// TODO: Получить доски только для авторизованного пользователя (по userID)
// TODO: Загружать доски из базы данных (ну или пока что просто из списка)
func (bs *BoardService) GetBoards() []models.Board {
	// Пока возвращаем пустой список
	return []models.Board{}
}

// AddBoard — создаёт новую доску, не обязательно, но как будто бы надо
// TODO: Проверить, что name не пустой
// TODO: Проверить, что пользователь авторизован (через session)
// TODO: Сохранить доску в БД с привязкой к userID
// TODO: Вернуть созданную доску
func (bs *BoardService) AddBoard(name string) (*models.Board, error) {
	// Пока просто возвращаем nil
	return nil, nil
}
