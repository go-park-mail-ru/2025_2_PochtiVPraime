package services

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

// BoardService — сервис для работы с досками
type BoardMemberService struct {
	//ListRepository  repository.ListsRepository
	//CardRepository  repository.CardsRepository
	BoardMemberRepository repository.BoardMemberRepository
	//UserRepository  repository.UserRepository
	// Здесь будут зависимости в будущем
}

// NewBoardService — конструктор (нужен для Dependency Injection) поботать эту тему ещё
func NewBoardMemberService(boardMemberRepository repository.BoardMemberRepository) *BoardMemberService {
	return &BoardMemberService{
		BoardMemberRepository: boardMemberRepository,
	}
}

// CreateBoardMember создает участника доски
func (bms *BoardMemberService) CreateBoardMember(ctx context.Context, boardMember *models.BoardMember, creatorId, boardId int64) (*models.BoardMember, error) {
	/*
		creator, err := bms.BoardMemberRepository.GetBoardMemberByUserId(ctx, boardId, creatorId)
		if creator.MemberRole != "owner" && creator.MemberRole != "admin" {
			return nil, fmt.Errorf("участников доски может добавлять только хозяин доски или админ: %w", err)
		}
	*/
	boardMember.BoardId = boardId
	boardMember, err := bms.BoardMemberRepository.CreateBoardMember(ctx, boardMember)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании участника доски: %w", err)
	}
	return boardMember, nil
}

// GetBoardMemberByUserId возвращает участника доски по userId и boardId
func (bms *BoardMemberService) GetBoardMemberByUserId(ctx context.Context, boardId, userId int64) (*models.BoardMember, error) {
	boardMember, err := bms.BoardMemberRepository.GetBoardMemberByUserId(ctx, boardId, userId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении участника доски по userId и boardId: %w", err)
	}
	return boardMember, nil
}

// GetMembersOfUser возвращает участников доски, которыми является наш пользователь
func (bms *BoardMemberService) GetMembersOfUser(ctx context.Context, userId int64) ([]*models.BoardMember, error) {
	var boardMembers, err = bms.BoardMemberRepository.GetBoardMembersByBoardId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return boardMembers, nil
}

// GetBoardMembers возвращает участников доски
func (bms *BoardMemberService) GetBoardMembers(ctx context.Context, boardId int64) ([]*models.BoardMember, error) {
	var boardMembers, err = bms.BoardMemberRepository.GetBoardMembersByBoardId(ctx, boardId)
	if err != nil {
		return nil, err
	}
	return boardMembers, nil
}

// ChangeRole возвращает участников доски
// Нужно будет сделать:
// Менять может только админ доски
// Может менять роли других пользователей
func (bms *BoardMemberService) ChangeRole(ctx context.Context, editorId, boardId int64, boardMember *models.BoardMember) error {
	/*
		editor, err := bms.BoardMemberRepository.GetBoardMemberByUserId(ctx, boardId, editorId)
		if err != nil {
			return fmt.Errorf("не удалось получить редактора: %w", err)
		}
		if editor.MemberRole != "owner" && editor.MemberRole != "admin" {
			return fmt.Errorf("удалять участников может только админ или хозяин доски : %w", err)
		}
	*/

	return bms.BoardMemberRepository.ChangeRole(ctx, boardMember.MemberRole, boardId, boardMember.UserId)
}

// DeleteBoardMembers позволяет только админам удалять участников доски
func (bms *BoardMemberService) DeleteBoardMember(ctx context.Context, editorId, boardId int64, boardMember *models.BoardMember) error {
	/*
		editor, err := bms.BoardMemberRepository.GetBoardMemberByUserId(ctx, boardId, editorId)
		if err != nil {
			return fmt.Errorf("не удалось получить редактора: %w", err)
		}
		if editor.MemberRole != "owner" && editor.MemberRole != "admin" {
			return fmt.Errorf("удалять участников может только админ : %w", err)
		}
		//var boardMemberId = bms.get
	*/
	boardMember.BoardId = boardId
	return bms.BoardMemberRepository.DeleteBoardMember(ctx, boardId, boardMember.UserId)
}
