package services

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

type ChecklistService struct {
	repo      repository.ChecklistRepository
	pointRepo repository.ChecklistPointRepository
}

func NewChecklistService(repo repository.ChecklistRepository, pointRepo repository.ChecklistPointRepository) *ChecklistService {
	return &ChecklistService{
		repo:      repo,
		pointRepo: pointRepo,
	}
}

// CreateChecklist создает новый чеклист с валидацией
func (s *ChecklistService) CreateChecklist(ctx context.Context, checklist *models.Checklist) error {
	if checklist.Title == "" {
		return fmt.Errorf("заголовок чеклиста не может быть пустым")
	}
	// Можно добавить дополнительные проверки или бизнес-логику
	return s.repo.CreateChecklist(ctx, checklist)
}

// GetChecklistByID возвращает чеклист по ID
func (s *ChecklistService) GetChecklistByID(ctx context.Context, id int64) (*models.ChecklistData, error) {
	// Получаем основной чеклист
	checklist, err := s.repo.GetChecklistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить чеклист: %w", err)
	}

	// Если чеклист не найден, возвращаем пустую структуру с cardID
	if checklist == nil {
		return &models.ChecklistData{
			CardId:          checklist.CardId,
			Title:           "Чеклист",
			ChecklistPoints: []models.ChecklistPoint{},
		}, nil
	}

	// Получаем все пункты чеклиста
	points, err := s.pointRepo.GetChecklistPointsByChecklistID(ctx, checklist.ID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пункты чеклиста: %w", err)
	}

	// Конвертируем []*models.ChecklistPoint в []models.ChecklistPoint
	checklistPoints := make([]models.ChecklistPoint, len(points))
	for i, point := range points {
		checklistPoints[i] = *point
	}

	// Собираем полные данные чеклиста
	checklistData := &models.ChecklistData{
		CardId:          checklist.CardId,
		Title:           checklist.Title,
		ChecklistPoints: checklistPoints,
		CreatedAt:       checklist.CreatedAt,
		UpdatedAt:       checklist.UpdatedAt,
	}

	return checklistData, nil
}

// GetChecklistsByCardID возвращает чеклисты по ID карточки
func (s *ChecklistService) GetChecklistsByCardID(ctx context.Context, cardID int64) ([]*models.ChecklistData, error) {
	// Получаем все чеклисты карточки
	checklists, err := s.repo.GetChecklistsByCardID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить чеклисты: %w", err)
	}

	// Если чеклистов нет, возвращаем пустой массив
	if len(checklists) == 0 {
		return []*models.ChecklistData{}, nil
	}

	var checklistDataList []*models.ChecklistData

	// Для каждого чеклиста получаем его пункты и формируем полные данные
	for _, checklist := range checklists {
		// Получаем пункты для текущего чеклиста
		points, err := s.pointRepo.GetChecklistPointsByChecklistID(ctx, checklist.ID)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить пункты чеклиста %d: %w", checklist.ID, err)
		}

		// Конвертируем []*models.ChecklistPoint в []models.ChecklistPoint
		checklistPoints := make([]models.ChecklistPoint, len(points))
		for i, point := range points {
			checklistPoints[i] = *point
		}

		// Собираем полные данные чеклиста
		checklistData := &models.ChecklistData{
			CardId:          checklist.CardId,
			Title:           checklist.Title,
			ChecklistPoints: checklistPoints,
			CreatedAt:       checklist.CreatedAt,
			UpdatedAt:       checklist.UpdatedAt,
		}

		checklistDataList = append(checklistDataList, checklistData)
	}

	return checklistDataList, nil
}

// UpdateChecklist обновляет чеклист после проверки
func (s *ChecklistService) UpdateChecklist(ctx context.Context, checklist *models.Checklist, checklistId int64) error {
	if checklist.Title == "" {
		return fmt.Errorf("заголовок чеклиста не может быть пустым")
	}
	checklist.ID = checklistId
	// Еще возможна проверка прав доступа или логика бизнес-процессов
	return s.repo.UpdateChecklist(ctx, checklist)
}

// DeleteChecklist удаляет чеклист по его ID
func (s *ChecklistService) DeleteChecklist(ctx context.Context, id int64) error {
	// Есть смысл проверить, существует ли чеклист перед удалением
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("чеклист не найден")
	}
	return s.repo.DeleteChecklist(ctx, id)
}

// UpdateTitle изменяет только название чеклиста
func (s *ChecklistService) UpdateTitle(ctx context.Context, id int64, title string) error {
	if title == "" {
		return fmt.Errorf("заголовок не может быть пустым")
	}
	return s.repo.UpdateTitle(ctx, id, title)
}
