package handlers

import (
	"book-cover-service/domains"
	"book-cover-service/handlers/tools"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPCoverLikeHandlers struct {
	coverLikeService services.CoverLikeService
	ctx              context.Context
}

func (lh *HTTPCoverLikeHandlers) HandleSetLike(w http.ResponseWriter, r *http.Request) {
	// Получить id пользователя и обложки
	var coverLike domains.CoverLike

	if err := json.NewDecoder(r.Body).Decode(&coverLike); err != nil {
		fmt.Println("ошибка получения данных лайка из тела запроса")
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Создать сущность лайка
	if err := lh.coverLikeService.SetLike(lh.ctx, coverLike.UserID, coverLike.CoverID); err != nil {
		fmt.Println("ошибка при записи данных лайка в таблицу")
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
}
