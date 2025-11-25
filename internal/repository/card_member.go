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

type СardMemberRepoImpl struct {
	db *sqlx.DB
}

func NewCardMemberRepository(db *sqlx.DB) CardMemberRepository {
	return &СardMemberRepoImpl{db: db}
}

// CreateCardMember создает новую связь между карточкой и участником доски
func (cmr *СardMemberRepoImpl) CreateCardMember(ctx context.Context, cardMember *models.CardMember) error {
	query := `
		INSERT INTO card_member (card_id, board_member_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	err := cmr.db.QueryRowContext(ctx, query,
		cardMember.CardID,
		cardMember.BoardMemberID,
		now,
		now,
	).Scan(&cardMember.ID)

	if err != nil {
		return fmt.Errorf("не удалось создать связь карточки и участника: %w", err)
	}

	cardMember.CreatedAt = now
	cardMember.UpdatedAt = now
	return nil
}

// GetCardMemberByID возвращает связь по ID
func (cmr *СardMemberRepoImpl) GetCardMemberByID(ctx context.Context, id int64) (*models.CardMember, error) {
	query := `
		SELECT id, card_id, board_member_id, created_at, updated_at
		FROM card_member
		WHERE id = $1
	`

	cardMember := &models.CardMember{}
	err := cmr.db.QueryRowContext(ctx, query, id).Scan(
		&cardMember.ID,
		&cardMember.CardID,
		&cardMember.BoardMemberID,
		&cardMember.CreatedAt,
		&cardMember.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("связь карточки и участника не найдена: %w", err)
		}
		return nil, fmt.Errorf("не удалось получить связь карточки и участника: %w", err)
	}

	return cardMember, nil
}

// GetCardMembersByCardID возвращает всех участников карточки
func (cmr *СardMemberRepoImpl) GetCardMembersByCardID(ctx context.Context, cardID int64) ([]*models.CardMember, error) {
	query := `
		SELECT id, card_id, board_member_id, created_at, updated_at
		FROM card_member
		WHERE card_id = $1
		ORDER BY created_at ASC
	`

	rows, err := cmr.db.QueryContext(ctx, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить участников карточки: %w", err)
	}
	defer rows.Close()

	var cardMembers []*models.CardMember
	for rows.Next() {
		cardMember := &models.CardMember{}
		err := rows.Scan(
			&cardMember.ID,
			&cardMember.CardID,
			&cardMember.BoardMemberID,
			&cardMember.CreatedAt,
			&cardMember.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные участника карточки: %w", err)
		}
		cardMembers = append(cardMembers, cardMember)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе участников карточки: %w", err)
	}

	return cardMembers, nil
}

// GetCardMembersByBoardMemberID возвращает все карточки участника
func (cmr *СardMemberRepoImpl) GetCardMembersByBoardMemberID(ctx context.Context, boardMemberID int64) ([]*models.CardMember, error) {
	query := `
		SELECT id, card_id, board_member_id, created_at, updated_at
		FROM card_member
		WHERE board_member_id = $1
		ORDER BY created_at ASC
	`

	rows, err := cmr.db.QueryContext(ctx, query, boardMemberID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить карточки участника: %w", err)
	}
	defer rows.Close()

	var cardMembers []*models.CardMember
	for rows.Next() {
		cardMember := &models.CardMember{}
		err := rows.Scan(
			&cardMember.ID,
			&cardMember.CardID,
			&cardMember.BoardMemberID,
			&cardMember.CreatedAt,
			&cardMember.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать данные карточки участника: %w", err)
		}
		cardMembers = append(cardMembers, cardMember)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе карточек участника: %w", err)
	}

	return cardMembers, nil
}

// DeleteCardMemberByCardAndBoardMember удаляет связь по card_id и board_member_id
func (cmr *СardMemberRepoImpl) DeleteCardMember(ctx context.Context, cardID, boardMemberID int64) error {
	query := `DELETE FROM card_member WHERE card_id = $1 AND board_member_id = $2`

	result, err := cmr.db.ExecContext(ctx, query, cardID, boardMemberID)
	if err != nil {
		return fmt.Errorf("не удалось удалить связь карточки и участника: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("связь карточки и участника не найдена")
	}

	return nil
}

// GetCardMember возвращает связь по card_id и board_member_id
func (cmr *СardMemberRepoImpl) GetCardMember(ctx context.Context, cardID, boardMemberID int64) (*models.CardMember, error) {
	query := `
		SELECT id, card_id, board_member_id, created_at, updated_at
		FROM card_member
		WHERE card_id = $1 AND board_member_id = $2
	`

	cardMember := &models.CardMember{}
	err := cmr.db.QueryRowContext(ctx, query, cardID, boardMemberID).Scan(
		&cardMember.ID,
		&cardMember.CardID,
		&cardMember.BoardMemberID,
		&cardMember.CreatedAt,
		&cardMember.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить связь карточки и участника: %w", err)
	}

	return cardMember, nil
}

// AddMemberToCard добавляет участника к карточке
func (r *СardMemberRepoImpl) AddMemberToCard(ctx context.Context, cardID int64, boardMemberID int64) error {
	query := `
        INSERT INTO card_member (card_id, board_member_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (card_id, board_member_id) DO NOTHING
    `

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, cardID, boardMemberID, now, now)
	if err != nil {
		return fmt.Errorf("не удалось добавить участника %d к карточке %d: %w", boardMemberID, cardID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return nil // пользователь уже добавлен к карточке
	}

	return nil
}

// RemoveMemberFromCard удаляет участника из карточки
func (r *СardMemberRepoImpl) RemoveMemberFromCard(ctx context.Context, cardID int64, boardMemberID int64) error {
	query := `DELETE FROM card_member WHERE card_id = $1 AND board_member_id = $2`

	result, err := r.db.ExecContext(ctx, query, cardID, boardMemberID)
	if err != nil {
		return fmt.Errorf("не удалось удалить участника %d из карточки %d: %w", boardMemberID, cardID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("участник %d не найден в карточке %d", boardMemberID, cardID)
	}

	return nil
}

// DeleteAllCardMembersByCardID удаляет всех участников карточки
func (r *СardMemberRepoImpl) DeleteAllCardMembersByCardID(ctx context.Context, cardID int64) error {
	query := `DELETE FROM card_member WHERE card_id = $1`

	_, err := r.db.ExecContext(ctx, query, cardID)
	if err != nil {
		return fmt.Errorf("не удалось удалить всех участников карточки: %w", err)
	}

	return nil
}

// DeleteAllCardMembersByBoardMemberID удаляет все связи участника доски
func (r *СardMemberRepoImpl) DeleteAllCardMembersByBoardMemberID(ctx context.Context, boardMemberID int64) error {
	query := `DELETE FROM card_member WHERE board_member_id = $1`

	_, err := r.db.ExecContext(ctx, query, boardMemberID)
	if err != nil {
		return fmt.Errorf("не удалось удалить все связи участника доски: %w", err)
	}

	return nil
}
