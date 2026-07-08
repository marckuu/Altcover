package services

import (
	"book-cover-service/core/domains"
	coverRepositoriesInterfaces "book-cover-service/internal/covers/repositories/interfaces"
	reactionRepositoriesInterfaces "book-cover-service/internal/reactions/repositories/interfaces"
	"context"

	"github.com/google/uuid"
)

type FavoritesService struct {
	FavoritesRepository reactionRepositoriesInterfaces.FavoritesRepository
	CoverRepository     coverRepositoriesInterfaces.CoverRepository
}

func NewFavoriteService(favoriteRepository reactionRepositoriesInterfaces.FavoritesRepository, coverRepository coverRepositoriesInterfaces.CoverRepository) *FavoritesService {
	return &FavoritesService{
		FavoritesRepository: favoriteRepository,
		CoverRepository:     coverRepository,
	}
}

func (f *FavoritesService) AddCoverToFavorites(ctx context.Context, userID uuid.UUID, coverID uuid.UUID) error {
	if err := f.FavoritesRepository.AddCoverToFavorites(ctx, userID, coverID); err != nil {
		return err
	}

	return nil
}

func (f *FavoritesService) GetFavoriteCovers(ctx context.Context, userID uuid.UUID, offset int, limit int) ([]domains.Cover, error) {
	// Получение ids избранных обложек пользователя с userID
	coversIDs, err := f.FavoritesRepository.GetFavoriteCoversIDs(ctx, userID, offset, limit)
	if err != nil {
		return []domains.Cover{}, err
	}

	// Получение самих обложек по этим ids
	covers, err := f.CoverRepository.GetCoversByIDs(ctx, coversIDs)
	if err != nil {
		return []domains.Cover{}, err
	}

	return covers, nil
}
