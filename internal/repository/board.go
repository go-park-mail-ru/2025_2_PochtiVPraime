package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/jmoiron/sqlx"
)

type BoardRepoImpl struct {
	DB *sqlx.DB
}

func NewBoardRepoImpl(db *sqlx.DB) BoardsRepository {
	return &BoardRepoImpl{
		DB: db,
	}
}

func (br *BoardRepoImpl) CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	query := `
		INSERT INTO board (owner_user_id, title, image_id, archived, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	var createdAt time.Time
	err := br.DB.QueryRowContext(ctx, query,
		board.OwnerId,
		board.Title,
		board.Image,
		board.Archived,
		time.Now(),
	).Scan(&board.ID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("не удалось создать доску: %w", err)
	}

	return board, nil
}

func (br *BoardRepoImpl) GetBoardById(ctx context.Context, id int64) (*models.Board, error) {
	query := `
		SELECT id, owner_user_id, title, archived, created_at
		FROM board
		WHERE id = $1
	`

	var board models.Board
	err := br.DB.QueryRowContext(ctx, query, id).Scan(
		&board.ID,
		&board.OwnerId,
		&board.Title,
		//&board.Image,
		&board.Archived,
		&board.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("доска не найдена")
		}
		return nil, fmt.Errorf("не удалось получить доску: %w", err)
	}

	return &board, nil
}

func (br *BoardRepoImpl) GetBoardsByOwner(ctx context.Context, ownerID int64) ([]*models.Board, error) {
	query := `
		SELECT id, owner_user_id, title, archived, created_at, updated_at
		FROM board
		WHERE owner_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := br.DB.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить доски: %w", err)
	}
	defer rows.Close()

	var boards []*models.Board
	for rows.Next() {
		var board models.Board
		err := rows.Scan(
			&board.ID,
			&board.OwnerId,
			&board.Title,
			//&board.Image,
			&board.Archived,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать доски: %w", err)
		}
		boards = append(boards, &board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по доскам: %w", err)
	}

	return boards, nil
}

func (br *BoardRepoImpl) UpdateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	query := `
		UPDATE board
		SET title = $1, archived = $2
		WHERE id = $3 AND owner_user_id = $4
		RETURNING created_at
	`

	var createdAt time.Time
	err := br.DB.QueryRowContext(ctx, query,
		board.Title,
		//board.Image,
		board.Archived,
		board.ID,
		board.OwnerId,
	).Scan(&createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("доска не найдена или нет доступа к ней")
		}
		return nil, fmt.Errorf("не удалось обновить доску: %w", err)
	}

	return board, nil
}

func (br *BoardRepoImpl) ArchiveBoard(ctx context.Context, id int64) error {
	query := `
		UPDATE board
		SET archived = true
		WHERE id = $1
	`

	result, err := br.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось заархивировать доску: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество заархивированных досок: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("доска не найдена или нет доступа к архивации")
	}

	return nil
}

func (br *BoardRepoImpl) RestoreBoard(ctx context.Context, id int64) error {
	query := `
		UPDATE board
		SET archived = false
		WHERE id = $1
	`

	result, err := br.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось восстановить доску: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество восстановленных досок: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("доска не найдена или нет доступа к восстановлению")
	}

	return nil
}

func (br *BoardRepoImpl) DeleteBoard(ctx context.Context, id int64) error {
	query := `DELETE FROM board WHERE id = $1`

	result, err := br.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить доску %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых досок: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("доска не найдена или нет доступа к удалению")
	}

	return nil
}
