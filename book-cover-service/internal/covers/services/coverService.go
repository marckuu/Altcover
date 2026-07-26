package services

import (
	"book-cover-service/core/domains"
	"book-cover-service/internal/covers/repositories/interfaces"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CoverService struct {
	coverRepository interfaces.CoverRepository
}

func NewCoverService(repository interfaces.CoverRepository) *CoverService {
	return &CoverService{
		coverRepository: repository,
	}
}

func (c *CoverService) GetCoversByUserID(ctx context.Context, offset int, limit int, userID uuid.UUID) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCoversByUserID(ctx, offset, limit, userID)
	if err != nil {
		return []domains.Cover{}, nil
	}

	return covers, nil
}

func (c *CoverService) GetCoverByID(ctx context.Context, coverID uuid.UUID) (domains.Cover, error) {
	cover, err := c.coverRepository.GetCoverByID(ctx, coverID)
	if err != nil {
		return domains.Cover{}, err
	}

	return cover, nil
}

func (c *CoverService) UpdateCover(ctx context.Context, coverID uuid.UUID, userID uuid.UUID, newCover domains.Cover) (domains.Cover, error) {
	cover, err := c.coverRepository.GetCoverByID(ctx, coverID)
	if err != nil {
		return domains.Cover{}, err
	}

	if cover.UserID != userID {
		return domains.Cover{}, fmt.Errorf("дизайнер не является автором данной обложки: %w", err)
	}

	savedCover, err := c.coverRepository.UpdateCover(ctx, newCover)
	if err != nil {
		return domains.Cover{}, err
	}

	return savedCover, nil
}

func (c *CoverService) AddCover(ctx context.Context, cover domains.Cover) (domains.Cover, error) {
	cover, err := c.coverRepository.AddCover(ctx, cover)
	if err != nil {
		return domains.Cover{}, err
	}
	return cover, nil
}

func (c *CoverService) GetCoversByIDs(ctx context.Context, coversIDs []uuid.UUID) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCoversByIDs(ctx, coversIDs)
	if err != nil {
		return []domains.Cover{}, err
	}

	return covers, nil
}

func (c *CoverService) GetMostLikedCovers(ctx context.Context, daysNumber int, offset int, limit int) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetMostLikedCoversForNDays(ctx, daysNumber, offset, limit)
	if err != nil {
		return []domains.Cover{}, err
	}

	return covers, nil
}

func (c *CoverService) GetCoversByBook(ctx context.Context, bookID uuid.UUID, offset int, limit int) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCoversByBook(ctx, bookID, offset, limit)
	if err != nil {
		return []domains.Cover{}, err
	}
	return covers, nil
}

func (c *CoverService) GetCovers(ctx context.Context, offset int, limit int) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCovers(ctx, offset, limit)
	if err != nil {
		return []domains.Cover{}, err
	}
	return covers, nil
}
