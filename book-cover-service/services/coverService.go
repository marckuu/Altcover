package services

import (
	"book-cover-service/db/repositories"
	"book-cover-service/domains"
	"context"
)

type CoverService struct {
	coverRepository repositories.CoverRepository
}

func (c *CoverService) GetCoversByDesignerID(ctx context.Context, offset int, limit int, designerID string) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCoversByDesignerID(ctx, offset, limit, designerID)
	if err != nil {
		return []domains.Cover{}, nil
	}

	return covers, nil
}

func (c *CoverService) GetCoversByUserID(ctx context.Context, offset int, limit int, userID string) ([]domains.Cover, error) {
	covers, err := c.coverRepository.GetCoversByUserID(ctx, offset, limit, userID)
	if err != nil {
		return []domains.Cover{}, nil
	}

	return covers, nil
}
