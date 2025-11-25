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

type commentRepo struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return &commentRepo{db: db}
}

// CreateComment создает новый комментарий
func (r *commentRepo) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comment (card_id, board_member_owner_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		comment.CardId,
		comment.BoardMemberOwnerId,
		comment.Content,
		now,
		now,
	).Scan(&comment.ID)

	if err != nil {
		return fmt.Errorf("не удалось создать комментарий: %w", err)
	}

	comment.CreatedAt = now
	comment.UpdatedAt = now
	return nil
}

// GetCommentByID возвращает комментарий по ID
func (r *commentRepo) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	query := `
		SELECT id, card_id, board_member_owner_id, content, created_at, updated_at
		FROM comment
		WHERE id = $1
	`

	comment := &models.Comment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.CardId,
		&comment.BoardMemberOwnerId,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("комментарий не найден: %w", err)
		}
		return nil, fmt.Errorf("не удалось получить комментарий: %w", err)
	}

	return comment, nil
}

// GetCommentsByCardID возвращает все комментарии карточки
func (r *commentRepo) GetCommentsByCardID(ctx context.Context, cardID int64) ([]*models.Comment, error) {
	query := `
		SELECT id, card_id, board_member_owner_id, content, created_at, updated_at
		FROM comment
		WHERE card_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить комментарии карточки: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.CardId,
			&comment.BoardMemberOwnerId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные комментария: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе комментариев: %w", err)
	}

	return comments, nil
}

// GetCommentsByBoardMemberID возвращает все комментарии участника доски
func (r *commentRepo) GetCommentsByBoardMemberID(ctx context.Context, boardMemberID int64) ([]*models.Comment, error) {
	query := `
		SELECT id, card_id, board_member_owner_id, content, created_at, updated_at
		FROM comment
		WHERE board_member_owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, boardMemberID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить комментарии участника: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.CardId,
			&comment.BoardMemberOwnerId,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные комментария: %w", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе комментариев: %w", err)
	}

	return comments, nil
}

// UpdateComment обновляет комментарий
func (r *commentRepo) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `
		UPDATE comment
		SET content = $1, updated_at = $2
		WHERE id = $3
	`

	comment.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		comment.Content,
		comment.UpdatedAt,
		comment.ID,
	)

	if err != nil {
		return fmt.Errorf("не удалось обновить комментарий: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("комментарий не найден")
	}

	return nil
}

// DeleteComment удаляет комментарий
func (r *commentRepo) DeleteComment(ctx context.Context, id int64) error {
	query := `DELETE FROM comment WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить комментарий: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("комментарий не найден")
	}

	return nil
}

// UpdateCommentContent обновляет только содержание комментария
func (r *commentRepo) UpdateCommentContent(ctx context.Context, id int64, content string) error {
	query := `
		UPDATE comment
		SET content = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, content, time.Now(), id)
	if err != nil {
		return fmt.Errorf("не удалось обновить содержание комментария: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("комментарий не найден")
	}

	return nil
}

// DeleteCommentsByCardID удаляет все комментарии карточки
func (r *commentRepo) DeleteCommentsByCardID(ctx context.Context, cardID int64) error {
	query := `DELETE FROM comment WHERE card_id = $1`

	_, err := r.db.ExecContext(ctx, query, cardID)
	if err != nil {
		return fmt.Errorf("не удалось удалить комментарии карточки: %w", err)
	}

	return nil
}

// DeleteCommentsByBoardMemberID удаляет все комментарии участника доски
func (r *commentRepo) DeleteCommentsByBoardMemberID(ctx context.Context, boardMemberID int64) error {
	query := `DELETE FROM comment WHERE board_member_owner_id = $1`

	_, err := r.db.ExecContext(ctx, query, boardMemberID)
	if err != nil {
		return fmt.Errorf("не удалось удалить комментарии участника: %w", err)
	}

	return nil
}

// GetCommentCountByCardID возвращает количество комментариев карточки
func (r *commentRepo) GetCommentCountByCardID(ctx context.Context, cardID int64) (int, error) {
	query := `SELECT COUNT(*) FROM comment WHERE card_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, cardID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить количество комментариев: %w", err)
	}

	return count, nil
}
