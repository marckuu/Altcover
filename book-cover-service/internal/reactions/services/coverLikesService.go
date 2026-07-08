package services

import (
	reactionRepositoriesInterfaces "book-cover-service/internal/reactions/repositories/interfaces"
	"context"
	"errors"

	"github.com/google/uuid"
)

var errLikeAlreadyExist = errors.New("лайк уже поставлен")

type CoverLikeService struct {
	coverLikeRepository reactionRepositoriesInterfaces.CoverLikeRepository
}

func NewCoverLikeService(repository reactionRepositoriesInterfaces.CoverLikeRepository) *CoverLikeService {
	return &CoverLikeService{
		coverLikeRepository: repository,
	}
}

func (ls *CoverLikeService) SetLike(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error {
	// Проверить есть ли уже лайк
	_, err := ls.coverLikeRepository.GetLike(ctx, userID, coverID)
	// Если нет ошибки значит лайк есть и его нельзя поставить снова
	if err == nil {
		return err
	}

	if err = ls.coverLikeRepository.AddLike(ctx, userID, coverID); err != nil {
		return err
	}

	return nil
}
