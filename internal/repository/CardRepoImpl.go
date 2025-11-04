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

type CardRepoImpl struct {
	DB *sqlx.DB
}

func NewCardRepoImpl(db *sqlx.DB) CardsRepository {
	return &CardRepoImpl{
		DB: db,
	}
}

// CreateCard создает новую карточку
func (cr *CardRepoImpl) CreateCard(ctx context.Context, card *models.Card) (*models.Card, error) {
	query := `
		INSERT INTO cards (
			author_board_member_id, 
			list_id, 
			content, 
			position, 
			complete_before,
			created_at,
			updated_at
		) VALUES (:author_board_member_id, :list_id, :content, :position, :complete_before, :created_at, :updated_at)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	rows, err := cr.DB.NamedQueryContext(ctx, query, card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan created card: %w", err)
		}
	}

	return card, nil
}

// GetCard возвращает карточку по ID
func (cr *CardRepoImpl) GetCard(ctx context.Context, id int64) (*models.Card, error) {
	query := `
		SELECT 
			id, 
			author_board_member_id, 
			list_id, 
			content, 
			position, 
			created_at, 
			updated_at, 
			complete_before
		FROM cards 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var card models.Card
	err := cr.DB.GetContext(ctx, &card, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("card not found")
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return &card, nil
}

// GetCardsByList возвращает все карточки в указанном списке
func (cr *CardRepoImpl) GetCardsByList(ctx context.Context, listID int64) ([]*models.Card, error) {
	query := `
		SELECT 
			id, 
			author_board_member_id, 
			list_id, 
			content, 
			position, 
			created_at, 
			updated_at, 
			complete_before
		FROM cards 
		WHERE list_id = $1 AND deleted_at IS NULL
		ORDER BY position ASC, created_at ASC
	`

	var cards []*models.Card
	err := cr.DB.SelectContext(ctx, &cards, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards by list: %w", err)
	}

	return cards, nil
}

// UpdateCard обновляет карточку
func (cr *CardRepoImpl) UpdateCard(ctx context.Context, card *models.Card) (*models.Card, error) {
	query := `
		UPDATE cards 
		SET 
			author_board_member_id = :author_board_member_id,
			list_id = :list_id,
			content = :content,
			position = :position,
			complete_before = :complete_before,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`

	card.UpdatedAt = time.Now()

	result, err := cr.DB.NamedExecContext(ctx, query, card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, errors.New("card not found")
	}

	return card, nil
}

// DeleteCard мягко удаляет карточку (устанавливает deleted_at)
func (cr *CardRepoImpl) DeleteCard(ctx context.Context, id int64) error {
	query := `
		UPDATE cards 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := cr.DB.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("card not found")
	}

	return nil
}

// UpdateCardPosition обновляет позицию карточки и/или перемещает в другой список
func (cr *CardRepoImpl) UpdateCardPosition(ctx context.Context, cardID int64, newPosition int, newListID int64) error {
	tx, err := cr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Проверяем существование карточки
	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM cards WHERE id = $1 AND deleted_at IS NULL)", cardID)
	if err != nil {
		return fmt.Errorf("failed to check card existence: %w", err)
	}
	if !exists {
		return errors.New("card not found")
	}

	// Обновляем позицию карточки
	query := `
		UPDATE cards 
		SET 
			position = $1,
			list_id = $2,
			updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := tx.ExecContext(ctx, query, newPosition, newListID, time.Now(), cardID)
	if err != nil {
		return fmt.Errorf("failed to update card position: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("card not found")
	}

	return tx.Commit()
}

// GetCardsByBoardMember возвращает все карточки созданные указанным участником доски
func (cr *CardRepoImpl) GetCardsByBoardMember(ctx context.Context, boardMemberID int64) ([]*models.Card, error) {
	query := `
		SELECT 
			id, 
			author_board_member_id, 
			list_id, 
			content, 
			position, 
			created_at, 
			updated_at, 
			complete_before
		FROM cards 
		WHERE author_board_member_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	var cards []*models.Card
	err := cr.DB.SelectContext(ctx, &cards, query, boardMemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards by board member: %w", err)
	}

	return cards, nil
}
