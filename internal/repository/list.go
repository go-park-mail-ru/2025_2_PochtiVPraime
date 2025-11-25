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

type ListRepoImpl struct {
	DB *sqlx.DB
}

func NewListRepoImpl(db *sqlx.DB) ListsRepository {
	return &ListRepoImpl{
		DB: db,
	}
}

// Сохраняет новый список
func (lr *ListRepoImpl) SaveList(ctx context.Context, list *models.List) (*models.List, error) {
	query := `
		INSERT INTO list (board_id, title, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	list.CreatedAt = now
	list.UpdatedAt = now

	err := lr.DB.QueryRowContext(
		ctx,
		query,
		list.BoardId,
		list.Title,
		list.Position,
		list.CreatedAt,
		list.UpdatedAt,
	).Scan(&list.ID, &list.CreatedAt, &list.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("не удалось создать список: %w", err)
	}

	return list, nil
}

// GetByID возвращает список по ID
func (lr *ListRepoImpl) GetListByID(ctx context.Context, id int64) (*models.List, error) {
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM list
		WHERE id = $1
	`

	var list models.List
	err := lr.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID,
		&list.BoardId,
		&list.Title,
		&list.Position,
		&list.CreatedAt,
		&list.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("список не найден")
		}
		return nil, fmt.Errorf("не удалось получить список по ИД: %w", err)
	}

	return &list, nil
}

// GetByBoardID возвращает все списки для указанной доски
func (lr *ListRepoImpl) GetListsByBoardID(ctx context.Context, boardID int64) ([]*models.List, error) {
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM list
		WHERE board_id = $1 
		ORDER BY position ASC, created_at ASC
	`

	rows, err := lr.DB.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить списки по ИД доски: %w", err)
	}
	defer rows.Close()

	var lists []*models.List
	for rows.Next() {
		var list models.List
		err := rows.Scan(
			&list.ID,
			&list.BoardId,
			&list.Title,
			&list.Position,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать списки: %w", err)
		}
		lists = append(lists, &list)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во времени иттерации по набору списков: %w", err)
	}

	return lists, nil
}

// Update обновляет существующий список
func (lr *ListRepoImpl) UpdateList(ctx context.Context, list *models.List) (*models.List, error) {
	query := `
		UPDATE list
		SET title = $1, position = $2, updated_at = $3
		WHERE id = $4 
		RETURNING updated_at
	`

	list.UpdatedAt = time.Now()

	err := lr.DB.QueryRowContext(
		ctx,
		query,
		list.Title,
		list.Position,
		list.UpdatedAt,
		list.ID,
	).Scan(&list.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("список не найден")
		}
		return nil, fmt.Errorf("не удалось обновить список: %w", err)
	}

	return list, nil
}

// Delete удаляет список
func (lr *ListRepoImpl) DeleteList(ctx context.Context, id int64) error {
	query := `DELETE FROM list WHERE id = $1`

	result, err := lr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить список: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых списков: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("список не найден")
	}

	return nil
}
