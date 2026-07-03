package services

import (
	"book-cover-service/internal/reactions/repositories"
	"context"
	"errors"

	"github.com/google/uuid"
)

var errLikeAlreadyExist = errors.New("лайк уже поставлен")

type CoverLikeService struct {
	coverLikeRepository repositories.CoverLikeRepository
}

func NewCoverLikeService(repository repositories.CoverLikeRepository) CoverLikeService {
	return CoverLikeService{
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
