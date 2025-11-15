package services

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

// SpportService — сервис для работы с досками
type SupportService struct {
	SupportRepository repository.SupportRepository
	// Здесь будут зависимости в будущем
}

// NewSupportService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewSupportService(supportRepository repository.SupportRepository) *SupportService {
	return &SupportService{
		SupportRepository: supportRepository,
	}
}
func (ss *SupportService) AddSupportForm(ctx context.Context, supportForm *models.SupportForm) error {

	/*
		if supportForm. == "" {
			return errors.New("Нет заголовка")
		}
		board.Archived = false
		board.Image = 1 //костыль, который потом уберём
	*/
	err := ss.SupportRepository.CreateSupportForm(ctx, supportForm)
	if err != nil {
		return err
	}
	return nil
}

// GetSupportForms — возвращает список форм
func (ss *SupportService) GetUserSupportForms(ctx context.Context, userId int64) ([]*models.SupportForm, error) {
	//var supportForms []*models.SupportForm
	var supportForms, err = ss.SupportRepository.GetSupportFormsByOwner(ctx, userId)
	if err != nil {
		return nil, err
	}

	/*
		for _, value := range rawUserBoards {
		}
	*/
	return supportForms, nil
}

// GetSupportForms — возвращает список форм
func (ss *SupportService) GetAllSupportForms(ctx context.Context) ([]*models.SupportForm, error) {
	//var supportForms []*models.SupportForm
	var supportForms, err = ss.SupportRepository.GetAllSupportForms(ctx)
	if err != nil {
		return nil, err
	}

	/*
		for _, value := range rawUserBoards {
		}
	*/
	return supportForms, nil
}

func (ss *SupportService) GetSupportFormById(ctx context.Context, formId int64) (*models.SupportForm, error) {
	// 1. Получаем базовую информацию о доске
	supportForm, err := ss.SupportRepository.GetSupportFormById(ctx, formId)
	if err != nil {
		log.Printf("Ошибка получения формы по ИД %d: %v", formId, err)
		return nil, err
	}
	return supportForm, nil
}

func (ss *SupportService) DeleteSupportForm(ctx context.Context, formId int64) error {
	err := ss.SupportRepository.DeleteSupportForm(ctx, formId)
	if err != nil {
		return err
	}
	return nil
}
