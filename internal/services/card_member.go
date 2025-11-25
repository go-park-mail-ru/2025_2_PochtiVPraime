package services

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

type CardMemberService struct {
	CardMemberRepository repository.CardMemberRepository
}

func NewCardMemberService(cardMemberRepository repository.CardMemberRepository) *CardMemberService {
	return &CardMemberService{
		CardMemberRepository: cardMemberRepository,
	}
}

func (s *CardMemberService) CreateCardMember(ctx context.Context, cardMember *models.CardMember, cardId int64) error {
	// Add validation or business logic here if needed
	cardMember.CardID = cardId
	return s.CardMemberRepository.CreateCardMember(ctx, cardMember)
}

func (s *CardMemberService) GetCardMemberByID(ctx context.Context, id int64) (*models.CardMember, error) {
	return s.CardMemberRepository.GetCardMemberByID(ctx, id)
}

func (s *CardMemberService) GetCardMembersByCardID(ctx context.Context, cardID int64) ([]*models.CardMember, error) {
	return s.CardMemberRepository.GetCardMembersByCardID(ctx, cardID)
}

func (s *CardMemberService) GetCardMembersByBoardMemberID(ctx context.Context, boardMemberID int64) ([]*models.CardMember, error) {
	return s.CardMemberRepository.GetCardMembersByBoardMemberID(ctx, boardMemberID)
}

func (s *CardMemberService) DeleteCardMember(ctx context.Context, cardID, boardMemberID int64) error {
	return s.CardMemberRepository.DeleteCardMember(ctx, cardID, boardMemberID)
}

func (s *CardMemberService) GetCardMember(ctx context.Context, cardID, boardMemberID int64) (*models.CardMember, error) {
	return s.CardMemberRepository.GetCardMember(ctx, cardID, boardMemberID)
}

func (s *CardMemberService) AddMemberToCard(ctx context.Context, cardID, boardMemberID int64) error {
	return s.CardMemberRepository.AddMemberToCard(ctx, cardID, boardMemberID)
}

func (s *CardMemberService) RemoveMemberFromCard(ctx context.Context, cardID, boardMemberID int64) error {
	return s.CardMemberRepository.RemoveMemberFromCard(ctx, cardID, boardMemberID)
}

func (s *CardMemberService) DeleteAllCardMembersByCardID(ctx context.Context, cardID int64) error {
	return s.CardMemberRepository.DeleteAllCardMembersByCardID(ctx, cardID)
}

func (s *CardMemberService) DeleteAllCardMembersByBoardMemberID(ctx context.Context, boardMemberID int64) error {
	return s.CardMemberRepository.DeleteAllCardMembersByBoardMemberID(ctx, boardMemberID)
}
