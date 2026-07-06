package transport

import (
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	"book-cover-service/core/tools"
	"book-cover-service/internal/reactions/services"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPCoverLikeHandlers struct {
	coverLikeService services.CoverLikeService
	ctx              context.Context
	logger           logs.Logger
}

func NewCoverLikeHandlers(coverLikeService services.CoverLikeService, ctx context.Context, logger logs.Logger) HTTPCoverLikeHandlers {
	return HTTPCoverLikeHandlers{
		coverLikeService: coverLikeService,
		ctx:              ctx,
		logger:           logger,
	}
}

func (l *HTTPCoverLikeHandlers) HandleSetLike(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		l.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coverIDRaw, ok := mux.Vars(r)["cover_id"]
	if !ok {
		l.logger.Error(fmt.Errorf("не удалось получить ID обложки из url: %w", err).Error())
		tools.SendErrorResponse(w, errors.New("не удалось получить ID обложки из url"), http.StatusBadRequest)
		return
	}

	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		l.logger.Error(fmt.Errorf("не удалось преобразовать id обложки uuid: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = l.coverLikeService.SetLike(l.ctx, userID, coverID); err != nil {
		l.logger.Error(fmt.Errorf("ошибка при записи данных лайка в таблицу: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
}
