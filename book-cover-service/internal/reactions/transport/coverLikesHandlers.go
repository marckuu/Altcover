package transport

import (
	"book-cover-service/core/errors"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	reactionServicesInterfaces "book-cover-service/internal/reactions/services/interfaces"
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPCoverLikeHandlers struct {
	coverLikeService reactionServicesInterfaces.CoverLikeService
	ctx              context.Context
	logger           logs.Logger
}

func NewCoverLikeHandlers(coverLikeService reactionServicesInterfaces.CoverLikeService, ctx context.Context, logger logs.Logger) HTTPCoverLikeHandlers {
	return HTTPCoverLikeHandlers{
		coverLikeService: coverLikeService,
		ctx:              ctx,
		logger:           logger,
	}
}

// @Summary Set like
// @Description Поставить лайк
// @Security ApiKeyAuth
// @Tags likes
// @Param cover_id path string true "Cover id"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /covers/{cover_id}/like [post]
func (l *HTTPCoverLikeHandlers) HandleSetLike(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		l.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("не удалось получить ID пользователя из токена: %w", err), http.StatusBadRequest)
		return
	}

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		l.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("ошибка преобразования строки в uuid: %w", err), http.StatusBadRequest)
		return
	}

	if err = l.coverLikeService.SetLike(l.ctx, userID, coverID); err != nil {
		l.logger.Error(fmt.Errorf("ошибка при записи данных лайка в таблицу: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("ошибка при записи данных лайка в таблицу: %w", err), http.StatusInternalServerError)
		return
	}
}
