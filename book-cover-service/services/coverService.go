package services

import (
	"book-cover-service/db/repositories"
	"book-cover-service/domains"
	"context"

	"github.com/google/uuid"
)

type CoverService struct {
	coverRepository repositories.CoverRepository
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

func (c *CoverService) UpdateCover(ctx context.Context, cover domains.Cover) error {
	if err := c.coverRepository.UpdateCover(ctx, cover); err != nil {
		return err
	}
	return nil
}

func (c *CoverService) AddCover(ctx context.Context, cover domains.Cover) error {
	if err := c.coverRepository.AddCover(ctx, cover); err != nil {
		return err
	}

	return nil
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
