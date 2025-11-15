package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/jmoiron/sqlx"
)

type SupportRepoImpl struct {
	DB *sqlx.DB
}

func NewSupportRepoImpl(db *sqlx.DB) SupportRepository {
	return &SupportRepoImpl{
		DB: db,
	}
}

func (sr *SupportRepoImpl) CreateSupportForm(ctx context.Context, supportForm *models.SupportForm) error {
	query := `
		INSERT INTO support_form (user_id, helper_id, form_type, form_status, text, contact_email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	var createdAt time.Time
	err := sr.DB.QueryRowContext(ctx, query,
		supportForm.UserId,
		supportForm.HelperId,
		supportForm.FormType,
		supportForm.FormStatus,
		supportForm.Text,
		supportForm.ContactEmail,
		time.Now(),
		time.Now(),
	).Scan(&supportForm.ID, &createdAt)

	if err != nil {
		return fmt.Errorf("не удалось создать обращение: %w", err)
	}

	return nil
}

func (sr *SupportRepoImpl) GetSupportFormById(ctx context.Context, id int64) (*models.SupportForm, error) {
	query := `
		SELECT id, user_id, helper_id, form_type, form_status, text, contact_email, created_at, updated_at
		FROM support_form
		WHERE id = $1
	`

	var supportForm models.SupportForm
	err := sr.DB.QueryRowContext(ctx, query, id).Scan(
		&supportForm.ID,
		&supportForm.UserId,
		&supportForm.HelperId,
		&supportForm.FormType,
		&supportForm.FormStatus,
		&supportForm.Text,
		&supportForm.ContactEmail,
		&supportForm.CreatedAt,
		&supportForm.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("форма не найдена")
		}
		return nil, fmt.Errorf("не удалось получить форму: %w", err)
	}

	return &supportForm, nil
}

func (sr *SupportRepoImpl) GetSupportFormsByOwner(ctx context.Context, ownerID int64) ([]*models.SupportForm, error) {
	query := `
		SELECT id, user_id, helper_id, form_type, form_status, text, contact_email, created_at, updated_at
		FROM support_form
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := sr.DB.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить формы пользователя: %w", err)
	}
	defer rows.Close()

	var supportForms []*models.SupportForm
	for rows.Next() {
		var supportForm models.SupportForm
		err := rows.Scan(
			&supportForm.ID,
			&supportForm.UserId,
			&supportForm.HelperId,
			&supportForm.FormType,
			&supportForm.FormStatus,
			&supportForm.Text,
			&supportForm.ContactEmail,
			&supportForm.CreatedAt,
			&supportForm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать формы: %w", err)
		}
		supportForms = append(supportForms, &supportForm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по формам: %w", err)
	}

	return supportForms, nil
}

func (sr *SupportRepoImpl) DeleteSupportForm(ctx context.Context, id int64) error {
	query := `DELETE FROM support_form WHERE id = $1`

	result, err := sr.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить форму %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых форм: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("форма не найдена или нет доступа к удалению")
	}

	return nil
}

func (sr *SupportRepoImpl) GetAllSupportForms(ctx context.Context) ([]*models.SupportForm, error) {
	query := `
		SELECT id, user_id, helper_id, form_type, form_status, text, contact_email, created_at, updated_at
		FROM support_form
		ORDER BY created_at DESC
	`

	rows, err := sr.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить формы пользователя: %w", err)
	}
	defer rows.Close()

	var supportForms []*models.SupportForm
	for rows.Next() {
		var supportForm models.SupportForm
		err := rows.Scan(
			&supportForm.ID,
			&supportForm.UserId,
			&supportForm.HelperId,
			&supportForm.FormType,
			&supportForm.FormStatus,
			&supportForm.Text,
			&supportForm.ContactEmail,
			&supportForm.CreatedAt,
			&supportForm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось считать доски: %w", err)
		}
		supportForms = append(supportForms, &supportForm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка во время иттерации по доскам: %w", err)
	}

	return supportForms, nil
}
