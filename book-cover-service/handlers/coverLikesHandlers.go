package handlers

import (
	"book-cover-service/handlers/tools"
	"book-cover-service/middleware"
	"book-cover-service/services"
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
}

func NewCoverLikeHandlers(coverLikeService services.CoverLikeService, ctx context.Context) HTTPCoverLikeHandlers {
	return HTTPCoverLikeHandlers{
		coverLikeService: coverLikeService,
		ctx:              ctx,
	}
}

func (lh *HTTPCoverLikeHandlers) HandleSetLike(w http.ResponseWriter, r *http.Request) {
	// Получить id пользователя и обложки

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coverIDRaw, ok := mux.Vars(r)["cover_id"]
	if !ok {
		fmt.Printf("не удалось получить ID обложки из url: %s", err)
		tools.SendErrorResponse(w, errors.New("не удалось получить ID обложки из url"), http.StatusBadRequest)
		return
	}

	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		fmt.Printf("не удалось преобразовать id обложки uuid: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создать сущность лайка
	if err = lh.coverLikeService.SetLike(lh.ctx, userID, coverID); err != nil {
		fmt.Println("ошибка при записи данных лайка в таблицу")
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
}
