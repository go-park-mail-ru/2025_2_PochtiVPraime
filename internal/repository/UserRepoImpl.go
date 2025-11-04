package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepoImpl struct {
	DB *sqlx.DB
}

func NewUserRepoImpl(db *sqlx.DB) UserRepository {
	return &UserRepoImpl{
		DB: db}
}

func (ur *UserRepoImpl) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO "user" (email, username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := ur.DB.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		[]byte(user.Password),
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
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
		return nil, err
	}

	return user, nil
}

// GetUserByID возвращает пользователя по ID
func (ur *UserRepoImpl) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, username, password, avatar_id, created_at, updated_at
		FROM "user"
		WHERE id = $1
	`

	user := &models.User{}
	err := ur.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.AvatarID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (ur *UserRepoImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, email, username, password, avatar_id, created_at, updated_at
		FROM "user"
		WHERE username = $1
	`

	user := &models.User{}
	err := ur.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.AvatarID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser обновляет данные пользователя
func (ur *UserRepoImpl) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		UPDATE "user" 
		SET email = $1, username = $2, password = $3, avatar_id = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`

	user.UpdatedAt = time.Now()

	err := ur.DB.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.Password,
		user.AvatarID,
		user.UpdatedAt,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
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

		return nil, err
	}

	return user, nil
}

// DeleteUser удаляет пользователя (soft delete)
func (ur *UserRepoImpl) DeleteUser(ctx context.Context, id int64) error {
	query := `
		DELETE FROM "user" 
		WHERE id = $1
	`

	result, err := ur.DB.ExecContext(ctx, query, id)
	if err != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}
