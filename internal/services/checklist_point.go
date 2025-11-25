package services

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

type ChecklistPointService struct {
	ChecklistPointRepository repository.ChecklistPointRepository
}

func NewChecklistPointService(checklistPointRepository repository.ChecklistPointRepository) *ChecklistPointService {
	return &ChecklistPointService{
		ChecklistPointRepository: checklistPointRepository,
	}
}

// CreateChecklistPoint создает новый пункт чеклиста с проверкой и логикой позиции
func (cps *ChecklistPointService) CreateChecklistPoint(ctx context.Context, point *models.ChecklistPoint, checklistId int64) error {

	if point.Position == 0 {
		maxPos, err := cps.ChecklistPointRepository.GetMaxPosition(ctx, point.ChecklistID)
		if err != nil {
			return fmt.Errorf("не удалось получить максимальную позицию: %w", err)
		}
		point.Position = maxPos + 1
	}
	point.ChecklistID = checklistId
	return cps.ChecklistPointRepository.CreateChecklistPoint(ctx, point)
}

func (cps *ChecklistPointService) GetChecklistPointByID(ctx context.Context, id int64) (*models.ChecklistPoint, error) {
	return cps.ChecklistPointRepository.GetChecklistPointByID(ctx, id)
}

func (cps *ChecklistPointService) GetChecklistPointsByChecklistID(ctx context.Context, checklistID int64) ([]*models.ChecklistPoint, error) {
	return cps.ChecklistPointRepository.GetChecklistPointsByChecklistID(ctx, checklistID)
}

func (cps *ChecklistPointService) UpdateChecklistPoint(ctx context.Context, point *models.ChecklistPoint) error {
	// Здесь может проводиться дополнительная валидация или проверка
	return cps.ChecklistPointRepository.UpdateChecklistPoint(ctx, point)
}

func (cps *ChecklistPointService) DeleteChecklistPoint(ctx context.Context, id int64) error {
	return cps.ChecklistPointRepository.DeleteChecklistPoint(ctx, id)
}

func (cps *ChecklistPointService) UpdateCheckedStatus(ctx context.Context, id int64, checked bool) error {
	return cps.ChecklistPointRepository.UpdateCheckedStatus(ctx, id, checked)
}

func (cps *ChecklistPointService) DeletePointsByChecklistID(ctx context.Context, checklistID int64) error {
	return cps.ChecklistPointRepository.DeletePointsByChecklistId(ctx, checklistID)
}
