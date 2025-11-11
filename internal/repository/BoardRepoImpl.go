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
		return nil, fmt.Errorf("failed to create board: %w", err)
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
			return nil, fmt.Errorf("board not found")
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
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
		return nil, fmt.Errorf("failed to get boards: %w", err)
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
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}
		boards = append(boards, &board)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating boards: %w", err)
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
			return nil, fmt.Errorf("board not found or access denied")
		}
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	return board, nil
}

func (r *BoardRepoImpl) ArchiveBoard(ctx context.Context, id int64) error {
	query := `
		UPDATE board
		SET archived = true
		WHERE id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to archive board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board not found or access denied")
	}

	return nil
}

func (r *BoardRepoImpl) RestoreBoard(ctx context.Context, id int64) error {
	query := `
		UPDATE board
		SET archived = false
		WHERE id = $1
	`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to restore board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board not found or access denied")
	}

	return nil
}

func (r *BoardRepoImpl) DeleteBoard(ctx context.Context, id int64) error {
	query := `DELETE FROM board WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board not found or access denied")
	}

	return nil
}
