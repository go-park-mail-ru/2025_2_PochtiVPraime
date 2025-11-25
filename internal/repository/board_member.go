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

type BoardMemberRepoImpl struct {
	DB *sqlx.DB
}

func NewBoardMemberRepoImpl(db *sqlx.DB) BoardMemberRepository {
	return &BoardMemberRepoImpl{
		DB: db}
}

func (bmr *BoardMemberRepoImpl) CreateBoardMember(ctx context.Context, boardMember *models.BoardMember) (*models.BoardMember, error) {
	query := `
		INSERT INTO board_member (user_id, board_id, member_role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	boardMember.CreatedAt = now
	boardMember.UpdatedAt = now

	err := bmr.DB.QueryRowContext(ctx, query,
		boardMember.UserId,
		boardMember.BoardId,
		boardMember.MemberRole,
		boardMember.CreatedAt,
		boardMember.UpdatedAt,
	).Scan(&boardMember.ID, &boardMember.CreatedAt, &boardMember.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return boardMember, nil
}

func (bmr *BoardMemberRepoImpl) GetBoardMembersByBoardId(ctx context.Context, boardID int64) ([]*models.BoardMember, error) {
	query := `
		SELECT id, user_id, board_id, member_role, created_at, updated_at
		FROM board_member
		WHERE board_id = $1
		ORDER BY created_at DESC
	`

	rows, err := bmr.DB.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить участников доски: %w", err)
	}
	defer rows.Close()

	var boardMembers []*models.BoardMember
	for rows.Next() {
		var boardMember models.BoardMember
		err := rows.Scan(
			&boardMember.ID,
			&boardMember.UserId,
			&boardMember.BoardId,
			&boardMember.MemberRole,
			&boardMember.CreatedAt,
			&boardMember.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать участников доски: %w", err)
		}
		boardMembers = append(boardMembers, &boardMember)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по участникам доски: %w", err)
	}

	return boardMembers, nil
}

// GetBoardMemberById возвращает участника доски по ID
func (bmr *BoardMemberRepoImpl) GetBoardMemberById(ctx context.Context, id int64) (*models.BoardMember, error) {
	query := `
		SELECT id, user_id, board_id, member_role, created_at, updated_at
		FROM board_member
		WHERE id = $1
	`

	boardMember := &models.BoardMember{}
	err := bmr.DB.QueryRowContext(ctx, query, id).Scan(
		boardMember.ID,
		boardMember.UserId,
		boardMember.BoardId,
		boardMember.MemberRole,
		boardMember.CreatedAt,
		boardMember.UpdatedAt,
	)

	//log.Println(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("board_member not found")
		}
		return nil, err
	}

	return boardMember, nil
}

// GetBoardMemberByUserId возвращает участника доски по ID пользователя
func (bmr *BoardMemberRepoImpl) GetBoardMemberByUserId(ctx context.Context, boardId, userId int64) (*models.BoardMember, error) {
	query := `
		SELECT id, user_id, board_id, member_role, created_at, updated_at
		FROM board_member
		WHERE user_id = $1 AND board_id = $2
	`

	boardMember := &models.BoardMember{}
	err := bmr.DB.QueryRowContext(ctx, query, userId, boardId).Scan(
		boardMember.ID,
		boardMember.UserId,
		boardMember.BoardId,
		boardMember.MemberRole,
		boardMember.CreatedAt,
		boardMember.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("board_member not found")
		}
		return nil, err
	}

	return boardMember, nil
}

// GetMembersOfUser возвращает участников, которым является пользователь
func (bmr *BoardMemberRepoImpl) GetMembersOfUser(ctx context.Context, userId int64) ([]*models.BoardMember, error) {
	query := `
		SELECT id, user_id, board_id, member_role, created_at, updated_at
		FROM board_member
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := bmr.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить участников: %w", err)
	}
	defer rows.Close()

	var boardMembers []*models.BoardMember
	for rows.Next() {
		var boardMember models.BoardMember
		err := rows.Scan(
			&boardMember.ID,
			&boardMember.UserId,
			&boardMember.BoardId,
			&boardMember.MemberRole,
			&boardMember.CreatedAt,
			&boardMember.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать участников пользователя: %w", err)
		}
		boardMembers = append(boardMembers, &boardMember)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по участникам пользователя: %w", err)
	}

	return boardMembers, nil
}

// UpdateUser обновляет данные пользователя
func (bmr *BoardMemberRepoImpl) ChangeRole(ctx context.Context, newRole string, boardId, userId int64) error {
	query := `
		UPDATE board_member 
		SET member_role = $1, updated_at = $2
		WHERE board_id = $3 AND user_id = $4
		RETURNING updated_at
	`

	updatedAt := time.Now()

	err := bmr.DB.QueryRowContext(ctx, query,
		newRole,
		updatedAt,
		boardId,
		userId,
	).Scan(&updatedAt)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser удаляет пользователя
func (bmr *BoardMemberRepoImpl) DeleteBoardMember(ctx context.Context, boardId, userId int64) error {
	query := `
		DELETE FROM board_member 
		WHERE board_id = $1 AND user_id = $2
	`

	result, err := bmr.DB.ExecContext(ctx, query, boardId, userId)
	if err != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return errors.New("board_member not found")
		}
		return err
	}
	return nil
}
