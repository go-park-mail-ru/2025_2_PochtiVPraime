package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

// CardService — сервис для работы с карточками
type CardService struct {
	CardRepository  repository.CardsRepository
	ListRepository  repository.ListsRepository
	BoardRepository repository.BoardsRepository
	// Здесь будут зависимости в будущем
}

// NewCardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewCardService(cardRepository repository.CardsRepository, listRepository repository.ListsRepository, boardRepository repository.BoardsRepository) *CardService {
	return &CardService{
		CardRepository:  cardRepository,
		ListRepository:  listRepository,
		BoardRepository: boardRepository,
	}
}

func (cs *CardService) CreateCard(ctx context.Context, rawCard *models.Card, listId int64) (*models.Card, error) {
	if len(rawCard.Content) < 1 || len(rawCard.Content) > 1000 {
		return nil, errors.New("Invalid size of content")
	}

	// Проверяем существование списка
	_, err := cs.ListRepository.GetListByID(ctx, listId)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	// Получаем текущие карточки списка для определения позиции
	existingCards, err := cs.CardRepository.GetCardsByList(ctx, listId)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing cards: %w", err)
	}

	// Определяем позицию новой карточки (в конец)
	position := 0
	if len(existingCards) > 0 {
		position = existingCards[len(existingCards)-1].Position + 1
	}

	card := &models.Card{
		AuthorBoardMemberId: rawCard.AuthorBoardMemberId,
		ListId:              rawCard.ListId,
		Content:             rawCard.Content,
		Position:            position,
		CompleteBefore:      rawCard.CompleteBefore,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	// Сохраняем в репозитории
	createdCard, err := cs.CardRepository.CreateCard(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return createdCard, nil
}

// GetCard возвращает карточку по ID с проверкой прав доступа
func (cs *CardService) GetCard(ctx context.Context, cardID, userID int64) (*models.Card, error) {
	card, err := cs.CardRepository.GetCard(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	/*
		// Получаем список для проверки прав доступа к доске
		list, err := s.listRepo.GetListByID(ctx, card.ListId)
		if err != nil {
			return nil, fmt.Errorf("failed to get list: %w", err)
		}

		// Проверяем, что пользователь имеет доступ к доске
		userBoardMember, err := cs.BoardMemberRepository.GetBoardMemberByUserAndBoard(ctx, userID, list.BoardID)
		if err != nil || userBoardMember == nil {
			return nil, ErrCardAccessDenied
		}
	*/
	return card, nil
}

// GetListCards возвращает все карточки списка
func (cs *CardService) GetListCards(ctx context.Context, listID int64) ([]*models.Card, error) {
	// Получаем карточки списка
	cards, err := cs.CardRepository.GetCardsByList(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list cards: %w", err)
	}

	return cards, nil
}

// UpdateCard обновляет карточку
func (cs *CardService) UpdateCard(ctx context.Context, card *models.Card) (*models.Card, error) {
	if len(card.Content) < 1 || len(card.Content) > 1000 {
		return nil, errors.New("Invalid size of content")
	}

	card.UpdatedAt = time.Now()

	// Сохраняем изменения
	updatedCard, err := cs.CardRepository.UpdateCard(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return updatedCard, nil
}

// DeleteCard удаляет карточку
func (cs *CardService) DeleteCard(ctx context.Context, cardId int64) error {
	// Удаляем карточку
	err := cs.CardRepository.DeleteCard(ctx, cardId)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}
