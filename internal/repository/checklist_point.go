package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/jmoiron/sqlx"
)

type ChecklistPointRepoImpl struct {
	db *sqlx.DB
}

func NewChecklistPointRepository(db *sqlx.DB) ChecklistPointRepository {
	return &ChecklistPointRepoImpl{db: db}
}

func (cpr *ChecklistPointRepoImpl) CreateChecklistPoint(ctx context.Context, point *models.ChecklistPoint) error {
	query := `
		INSERT INTO checklist_point (checklist_id, content, checked, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	err := cpr.db.QueryRowContext(ctx, query,
		point.ChecklistID,
		point.Content,
		point.Checked,
		point.Position,
		now,
		now,
	).Scan(&point.ID)

	if err != nil {
		return fmt.Errorf("failed to create checklist point: %w", err)
	}

	point.CreatedAt = now
	point.UpdatedAt = now
	return nil
}

// GetChecklistPointByID возвращает пункт чеклиста по ID
func (cpr *ChecklistPointRepoImpl) GetChecklistPointByID(ctx context.Context, id int64) (*models.ChecklistPoint, error) {
	query := `
		SELECT id, checklist_id, content, checked, position, created_at, updated_at
		FROM checklist_point
		WHERE id = $1
	`

	point := &models.ChecklistPoint{}
	err := cpr.db.QueryRowContext(ctx, query, id).Scan(
		&point.ID,
		&point.ChecklistID,
		&point.Content,
		&point.Checked,
		&point.Position,
		&point.CreatedAt,
		&point.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пункт чеклиста не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении пункта челиста: %w", err)
	}

	return point, nil
}

// GetChecklistPointsByChecklistID возвращает все пункты чеклиста
func (cpr *ChecklistPointRepoImpl) GetChecklistPointsByChecklistID(ctx context.Context, checklistID int64) ([]*models.ChecklistPoint, error) {
	query := `
		SELECT id, checklist_id, content, checked, position, created_at, updated_at
		FROM checklist_point
		WHERE checklist_id = $1
		ORDER BY position ASC, created_at ASC
	`

	rows, err := cpr.db.QueryContext(ctx, query, checklistID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пункта челиста: %w", err)
	}
	defer rows.Close()

	var points []*models.ChecklistPoint
	for rows.Next() {
		point := &models.ChecklistPoint{}
		err := rows.Scan(
			&point.ID,
			&point.ChecklistID,
			&point.Content,
			&point.Checked,
			&point.Position,
			&point.CreatedAt,
			&point.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать пункт чеклиста: %w", err)
		}
		points = append(points, point)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при иттерации по пунктам чеклиста: %w", err)
	}

	return points, nil
}

// UpdateChecklistPoint обновляет пункт чеклиста
func (cpr *ChecklistPointRepoImpl) UpdateChecklistPoint(ctx context.Context, point *models.ChecklistPoint) error {
	query := `
		UPDATE checklist_point
		SET content = $1, checked = $2, position = $3, updated_at = $4
		WHERE id = $5
	`

	point.UpdatedAt = time.Now()
	result, err := cpr.db.ExecContext(ctx, query,
		point.Content,
		point.Checked,
		point.Position,
		point.UpdatedAt,
		point.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить пункт чеклиста: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновленных пунктов чеклиста: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пункт чеклиста не найден")
	}

	return nil
}

// DeleteChecklistPoint удаляет пункт чеклиста
func (cpr *ChecklistPointRepoImpl) DeleteChecklistPoint(ctx context.Context, id int64) error {
	query := `DELETE FROM checklist_point WHERE id = $1`

	result, err := cpr.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить пункт чеклиста: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удаленных пунктов чеклиста: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пункт чеклиста не найден")
	}

	return nil
}

// UpdateCheckedStatus обновляет только статус checked
func (cpr *ChecklistPointRepoImpl) UpdateCheckedStatus(ctx context.Context, id int64, checked bool) error {
	query := `
		UPDATE checklist_point
		SET checked = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := cpr.db.ExecContext(ctx, query, checked, time.Now(), id)
	if err != nil {
		return fmt.Errorf("не удалось изменить статус пункта чеклиста: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество изменных статусов пунктов чеклистов: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пункт чеклиста не найден")
	}

	return nil
}

// DeletePointsByChecklistID удаляет все пункты чеклиста
func (cpr *ChecklistPointRepoImpl) DeletePointsByChecklistId(ctx context.Context, checklistID int64) error {
	query := `DELETE FROM checklist_points WHERE checklist_id = $1`

	_, err := cpr.db.ExecContext(ctx, query, checklistID)
	if err != nil {
		return fmt.Errorf("не удалось удалить пункты чеклиста: %w", err)
	}

	return nil
}

func (cpr *ChecklistPointRepoImpl) GetMaxPosition(ctx context.Context, checklistID int64) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), 0)
		FROM checklist_point
		WHERE checklist_id = $1
	`

	var maxPosition int
	err := cpr.db.QueryRowContext(ctx, query, checklistID).Scan(&maxPosition)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить максимальную позицию пункта чеклиста: %w", err)
	}

	return maxPosition, nil
}
