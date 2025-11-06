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
		return nil, fmt.Errorf("failed to create list: %w", err)
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
			return nil, errors.New("List not found")
		}
		return nil, fmt.Errorf("failed to get list by id: %w", err)
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
		return nil, fmt.Errorf("failed to get lists by board id: %w", err)
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
			return nil, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, &list)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lists: %w", err)
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
			return nil, errors.New("List not found")
		}
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return list, nil
}

// Delete удаляет список (мягкое удаление)
func (lr *ListRepoImpl) DeleteList(ctx context.Context, id int64) error {
	query := `DELETE FROM board WHERE id = $1`

	result, err := lr.DB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("List not found")
	}

	return nil
}
