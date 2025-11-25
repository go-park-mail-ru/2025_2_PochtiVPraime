package services

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/models"
	"github.com/go-park-mail-ru/2025_2_PochtiVPraime/internal/repository"
)

type CommentService struct {
	CommentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepository: commentRepository,
	}
}

// CreateComment создает новый комментарий с проверкой и логикой
func (cs *CommentService) CreateComment(ctx context.Context, comment *models.Comment, boardMemberID int64) error {
	// Проверяем, что контент не пустой
	if comment.Content == "" {
		return fmt.Errorf("контент комментария не может быть пустым")
	}

	// Проверяем длину контента
	if len(comment.Content) > 5000 {
		return fmt.Errorf("контент комментария слишком длинный")
	}

	// Устанавливаем автора комментария
	comment.BoardMemberOwnerId = boardMemberID

	return cs.CommentRepository.CreateComment(ctx, comment)
}

// GetCommentByID возвращает комментарий по ID
func (cs *CommentService) GetCommentByID(ctx context.Context, id int64) (*models.Comment, error) {
	return cs.CommentRepository.GetCommentByID(ctx, id)
}

// GetCommentsByCardID возвращает все комментарии карточки
func (cs *CommentService) GetCommentsByCardID(ctx context.Context, cardID int64) ([]*models.Comment, error) {
	return cs.CommentRepository.GetCommentsByCardID(ctx, cardID)
}

// GetCommentsByBoardMemberID возвращает все комментарии участника доски
func (cs *CommentService) GetCommentsByBoardMemberID(ctx context.Context, boardMemberID int64) ([]*models.Comment, error) {
	return cs.CommentRepository.GetCommentsByBoardMemberID(ctx, boardMemberID)
}

// UpdateComment обновляет комментарий
func (cs *CommentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	// Проверяем, что контент не пустой
	if comment.Content == "" {
		return fmt.Errorf("контент комментария не может быть пустым")
	}

	// Проверяем длину контента
	if len(comment.Content) > 5000 {
		return fmt.Errorf("контент комментария слишком длинный")
	}

	return cs.CommentRepository.UpdateComment(ctx, comment)
}

// DeleteComment удаляет комментарий
func (cs *CommentService) DeleteComment(ctx context.Context, id int64) error {
	return cs.CommentRepository.DeleteComment(ctx, id)
}

// UpdateCommentContent обновляет только содержание комментария
func (cs *CommentService) UpdateCommentContent(ctx context.Context, id int64, content string) error {
	// Проверяем, что контент не пустой
	if content == "" {
		return fmt.Errorf("контент комментария не может быть пустым")
	}

	// Проверяем длину контента
	if len(content) > 5000 {
		return fmt.Errorf("контент комментария слишком длинный")
	}

	return cs.CommentRepository.UpdateCommentContent(ctx, id, content)
}

// DeleteCommentsByCardID удаляет все комментарии карточки
func (cs *CommentService) DeleteCommentsByCardID(ctx context.Context, cardID int64) error {
	return cs.CommentRepository.DeleteCommentsByCardID(ctx, cardID)
}

// DeleteCommentsByBoardMemberID удаляет все комментарии участника доски
func (cs *CommentService) DeleteCommentsByBoardMemberID(ctx context.Context, boardMemberID int64) error {
	return cs.CommentRepository.DeleteCommentsByBoardMemberID(ctx, boardMemberID)
}

// GetCommentCountByCardID возвращает количество комментариев карточки
func (cs *CommentService) GetCommentCountByCardID(ctx context.Context, cardID int64) (int, error) {
	return cs.CommentRepository.GetCommentCountByCardID(ctx, cardID)
}

// CommentExists проверяет существование комментария
func (cs *CommentService) CommentExists(ctx context.Context, id int64) (bool, error) {
	comment, err := cs.CommentRepository.GetCommentByID(ctx, id)
	if err != nil {
		return false, err
	}
	return comment != nil, nil
}

// CanUserEditComment проверяет, может ли пользователь редактировать комментарий
func (cs *CommentService) CanUserEditComment(ctx context.Context, commentID, boardMemberID int64) (bool, error) {
	comment, err := cs.CommentRepository.GetCommentByID(ctx, commentID)
	if err != nil {
		return false, err
	}

	// Пользователь может редактировать только свои комментарии
	return comment.BoardMemberOwnerId == boardMemberID, nil
}

// CanUserDeleteComment проверяет, может ли пользователь удалить комментарий
func (cs *CommentService) CanUserDeleteComment(ctx context.Context, commentID, boardMemberID int64) (bool, error) {
	comment, err := cs.CommentRepository.GetCommentByID(ctx, commentID)
	if err != nil {
		return false, err
	}

	// Пользователь может удалять только свои комментарии
	return comment.BoardMemberOwnerId == boardMemberID, nil
}
