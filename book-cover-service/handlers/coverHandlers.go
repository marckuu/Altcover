package handlers

import (
	"book-cover-service/dto"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPCoverHandlers struct {
	coverService services.CoverService
}

func SendErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	errorResponse := dto.NewErrorResponse(err.Error())
	if convertErr := json.NewEncoder(w).Encode(errorResponse); convertErr != nil {
		fmt.Println("Ошибка при записи ответа с информацией об ошибке")
	}
}

func (ch *HTTPCoverHandlers) HandleGetCoversByDesignerID(w http.ResponseWriter, r *http.Request, ctx context.Context, offset int, limit int) {
	designerID := mux.Vars(r)["designer_id"]

	covers, err := ch.coverService.GetCoversByDesignerID(ctx, offset, limit, designerID)
	if err != nil {
		fmt.Printf("Ошибка при получении списка обложек дизайнера: %s", err)
		SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("Ошибка при записи ответа с полученными обложками")
		// Логировать оишбку
	}
}

func (ch *HTTPCoverHandlers) HandleGetCoverByUserID(w http.ResponseWriter, r *http.Request, ctx context.Context, offset int, limit int) {
	userID := mux.Vars(r)["user_id"]

	// Получить список обложек по полученным ID

	covers, err := ch.coverService.GetCoversByUserID(ctx, offset, limit, userID)
	if err != nil {
		fmt.Printf("Ошибка при получении списка обложек дизайнера: %s", err)
		SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("Ошибка при записи ответа с полученными обложками")
		// Логировать оишбку
	}
}
