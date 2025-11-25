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

type ChecklistRepoImpl struct {
	db *sqlx.DB
}

func NewChecklistRepository(db *sqlx.DB) ChecklistRepository {
	return &ChecklistRepoImpl{
		db: db,
	}
}

// CreateChecklist создает новый чеклист
func (clr *ChecklistRepoImpl) CreateChecklist(ctx context.Context, checklist *models.Checklist) error {
	query := `
		INSERT INTO checklist (card_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err := clr.db.QueryRowContext(ctx, query,
		checklist.CardId,
		checklist.Title,
		now,
		now,
	).Scan(&checklist.ID)

	if err != nil {
		return fmt.Errorf("не удалось создать чеклист: %w", err)
	}

	checklist.CreatedAt = now
	checklist.UpdatedAt = now
	return nil
}

// GetChecklistByID возвращает чеклист по ID
func (clr *ChecklistRepoImpl) GetChecklistByID(ctx context.Context, id int64) (*models.Checklist, error) {
	query := `
		SELECT id, card_id, title, created_at, updated_at
		FROM checklist
		WHERE id = $1
	`

	checklist := &models.Checklist{}
	err := clr.db.QueryRowContext(ctx, query, id).Scan(
		&checklist.ID,
		&checklist.CardId,
		&checklist.Title,
		&checklist.CreatedAt,
		&checklist.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("чеклист не найден: %w", err)
		}
		return nil, fmt.Errorf("не далось получить чеклист: %w", err)
	}

	return checklist, nil
}

// GetChecklistsByCardID возвращает чеклисты по ID карточки
func (clr *ChecklistRepoImpl) GetChecklistsByCardID(ctx context.Context, cardID int64) ([]*models.Checklist, error) {
	query := `
        SELECT id, card_id, title, created_at, updated_at
        FROM checklist
        WHERE card_id = $1
        ORDER BY created_at ASC
    `

	rows, err := clr.db.QueryContext(ctx, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить чеклисты по карточке: %w", err)
	}
	defer rows.Close()

	var checklists []*models.Checklist
	for rows.Next() {
		checklist := &models.Checklist{}
		err := rows.Scan(
			&checklist.ID,
			&checklist.CardId,
			&checklist.Title,
			&checklist.CreatedAt,
			&checklist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные чеклиста: %w", err)
		}
		checklists = append(checklists, checklist)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе чеклистов: %w", err)
	}

	return checklists, nil
}

// UpdateChecklist обновляет чеклист
func (clr *ChecklistRepoImpl) UpdateChecklist(ctx context.Context, checklist *models.Checklist) error {
	query := `
		UPDATE checklist
		SET title = $1, updated_at = $2
		WHERE id = $3
	`

	checklist.UpdatedAt = time.Now()
	result, err := clr.db.ExecContext(ctx, query,
		checklist.Title,
		checklist.UpdatedAt,
		checklist.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить чеклист: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных чеклистов: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("чеклист не найден")
	}

	return nil
}

// DeleteChecklist удаляет чеклист
func (clr *ChecklistRepoImpl) DeleteChecklist(ctx context.Context, id int64) error {
	query := `DELETE FROM checklist WHERE id = $1`

	result, err := clr.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить чеклист: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых чеклистов: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("чеклист не найден")
	}

	return nil
}

// UpdateTitle обновляет только заголовок чеклиста
func (clr *ChecklistRepoImpl) UpdateTitle(ctx context.Context, id int64, title string) error {
	query := `
		UPDATE checklist
		SET title = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := clr.db.ExecContext(ctx, query, title, time.Now(), id)
	if err != nil {
		return fmt.Errorf("не удалось обновить заголовок чеклиста: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество чеклистов с обновленным заголовком: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("чеклист не найден")
	}

	return nil
}

// Exists проверяет существование чеклиста
func (clr *ChecklistRepoImpl) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM checklist WHERE id = $1)`

	var exists bool
	err := clr.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить существование чеклиста: %w", err)
	}

	return exists, nil
}
