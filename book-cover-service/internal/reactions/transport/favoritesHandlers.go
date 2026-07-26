package transport

import (
	"book-cover-service/core/errors"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	tools2 "book-cover-service/core/tools"
	reactionServicesInterfaces "book-cover-service/internal/reactions/services/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPFavoritesHandlers struct {
	favoritesService reactionServicesInterfaces.FavoritesService
	ctx              context.Context
	logger           logs.Logger
}

func NewHTTPFavoritesHandlers(favoriteService reactionServicesInterfaces.FavoritesService, ctx context.Context, logger logs.Logger) HTTPFavoritesHandlers {
	return HTTPFavoritesHandlers{
		favoritesService: favoriteService,
		ctx:              ctx,
		logger:           logger,
	}
}

// @Summary Add cover to favorites
// @Description Добавить обложку в избранное
// @Security ApiKeyAuth
// @Tags favorites
// @Param cover_id path string true "Cover id"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /covers/{cover_id}/favorite [post]
func (f *HTTPFavoritesHandlers) HandleAddCoverToFavorites(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		f.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("не удалось получить ID пользователя из токена: %w", err), http.StatusBadRequest)
		return
	}

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		f.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("ошибка преобразования строки в uuid: %w", err), http.StatusInternalServerError)
		return
	}

	if err = f.favoritesService.AddCoverToFavorites(f.ctx, userID, coverID); err != nil {
		f.logger.Error(fmt.Errorf("ошибка добавления обложки в избранное: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("ошибка добавления обложки в избранное: %w", err), http.StatusInternalServerError)
		return
	}
}

// @Summary Get my favorite covers
// @Description Получить обложки из избранного текущего пользователя
// @Security ApiKeyAuth
// @Tags favorites
// @Param offset query string true "Pagination offset"
// @Param limit query string true "Pagination limit"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /covers/me/favorites [get]
func (f *HTTPFavoritesHandlers) HandleGetMyFavoriteCovers(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools2.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		f.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		f.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("не удалось получить ID пользователя из токена: %w", err), http.StatusBadRequest)
		return
	}

	covers, err := f.favoritesService.GetFavoriteCovers(f.ctx, userID, int(offset), int(limit))
	if err != nil {
		f.logger.Error(fmt.Errorf("ошибка при получении избранных обложек: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("ошибка при получении избранных обложек: %w", err), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		f.logger.Error(fmt.Errorf("ошибка при получении избранных обложек: %w", err).Error())
	}
}
