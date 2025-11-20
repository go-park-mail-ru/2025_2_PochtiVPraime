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
		INSERT INTO card (
			author_board_member_id, 
			list_id, 
			content, 
			position, 
			complete_before,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	// Выполняем запрос с позиционными параметрами
	err := cr.DB.QueryRowContext(
		ctx,
		query,
		card.AuthorBoardMemberId,
		card.ListId,
		card.Content,
		card.Position,
		card.CompleteBefore,
		card.CreatedAt,
		card.UpdatedAt,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("не удалось создать карточку: %w", err)
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
			completed, 
			created_at, 
			updated_at, 
			complete_before
		FROM card 
		WHERE id = $1
	`

	var card models.Card
	err := cr.DB.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.AuthorBoardMemberId,
		&card.ListId,
		&card.Content,
		&card.Position,
		&card.Completed,
		&card.CreatedAt,
		&card.UpdatedAt,
		&card.CompleteBefore,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("карточка не найдена: %w", err)
		}
		return nil, fmt.Errorf("не удалось найти карточку: %w", err)
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
			completed,
			created_at, 
			updated_at, 
			complete_before
		FROM card 
		WHERE list_id = $1
		ORDER BY position ASC, created_at ASC
	`

	rows, err := cr.DB.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти карточку по списку: %w", err)
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		var card models.Card
		err := rows.Scan(
			&card.ID,
			&card.AuthorBoardMemberId,
			&card.ListId,
			&card.Content,
			&card.Position,
			&card.Completed,
			&card.CreatedAt,
			&card.UpdatedAt,
			&card.CompleteBefore,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить карточки: %w", err)
		}
		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по карточкам: %w", err)
	}

	return cards, nil
}

func (cr *CardRepoImpl) UpdateCard(ctx context.Context, card *models.Card) (*models.Card, error) {
	query := `
        UPDATE card 
        SET 
            content = $1,
            completed = $2,
            updated_at = $3,
            complete_before = $4
        WHERE id = $5
        RETURNING updated_at
    `

	card.UpdatedAt = time.Now()

	var updatedAt time.Time
	err := cr.DB.QueryRowContext(ctx, query,
		card.Content,
		card.Completed,
		card.UpdatedAt,
		card.CompleteBefore,
		card.ID,
	).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("карточка не найдена")
		}
		return nil, fmt.Errorf("не удалось обновить карточку: %w", err)
	}

	card.UpdatedAt = updatedAt
	return card, nil
}

// Удаляет карточку
func (cr *CardRepoImpl) DeleteCard(ctx context.Context, id int64) error {
	query := `DELETE FROM card WHERE id = $1`

	result, err := cr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить карточку: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных карточек: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("карточка не неайдена")
	}

	return nil
}

// UpdateCardPosition обновляет позицию карточки и/или перемещает в другой список
func (cr *CardRepoImpl) UpdateCardPosition(ctx context.Context, cardID int64, newPosition int, newListID int64) error {
	tx, err := cr.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}
	defer tx.Rollback()

	// Проверяем существование карточки
	var exists bool
	err = tx.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM cards WHERE id = $1 AND deleted_at IS NULL)", cardID)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование карточки: %w", err)
	}
	if !exists {
		return errors.New("карточка не найдена")
	}

	// Обновляем позицию карточки
	query := `
		UPDATE card 
		SET 
			position = $1,
			list_id = $2,
			updated_at = $3
		WHERE id = $4
	`

	result, err := tx.ExecContext(ctx, query, newPosition, newListID, cardID)
	if err != nil {
		return fmt.Errorf("неудалось обновить положение карточки: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных карточек: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("карточка не найдена")
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
			completed, 
			created_at, 
			updated_at, 
			complete_before
		FROM card 
		WHERE author_board_member_id = $1
		ORDER BY created_at DESC
	`

	var cards []*models.Card
	err := cr.DB.SelectContext(ctx, &cards, query, boardMemberID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить карточки созданные участником доски: %w", err)
	}

	return cards, nil
}
