package repository

import (
	"context"
	"database/sql"
	"errors"
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
		/*
			// Обработка ошибок уникальности
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					if pqErr.Constraint == "user_email_key" {
						return nil, errors.New("email already exists")
					}
					if pqErr.Constraint == "user_username_key" {
						return nil, errors.New("username already exists")
					}
				}
			}
		*/
		return nil, err
	}

	return boardMember, nil
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
func (bmr *BoardMemberRepoImpl) GetBoardMemberByUserId(ctx context.Context, userId int64) (*models.BoardMember, error) {
	query := `
		SELECT id, user_id, board_id, member_role, created_at, updated_at
		FROM board_member
		WHERE user_id = $1
	`

	boardMember := &models.BoardMember{}
	err := bmr.DB.QueryRowContext(ctx, query, userId).Scan(
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

// UpdateUser обновляет данные пользователя
func (bmr *BoardMemberRepoImpl) ChangeRole(ctx context.Context, newRole string, memberId int64) error {
	query := `
		UPDATE board_member 
		SET member_role = $1, update_at = $2
		WHERE id = $3
		RETURNING updated_at
	`

	updatedAt := time.Now()

	err := bmr.DB.QueryRowContext(ctx, query,
		newRole,
		updatedAt,
		memberId,
	).Scan(&updatedAt)

	if err != nil {
		/*
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New("user not found")
			}

			// Обработка ошибок уникальности
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					if pqErr.Constraint == "users_email_key" {
						return nil, errors.New("email already exists")
					}
					if pqErr.Constraint == "users_username_key" {
						return nil, errors.New("username already exists")
					}
				}
			}
		*/

		return err
	}

	return nil
}

// DeleteUser удаляет пользователя (soft delete)
func (bmr *BoardMemberRepoImpl) DeleteBoardMember(ctx context.Context, id int64) error {
	query := `
		DELETE FROM board_member 
		WHERE id = $1
	`

	result, err := bmr.DB.ExecContext(ctx, query, id)
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
