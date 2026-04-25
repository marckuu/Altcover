package handlers

import (
	"book-cover-service/handlers/tools"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPFavoritesHandlers struct {
	favoritesService services.FavoritesService
	ctx              context.Context
}

func (fh *HTTPFavoritesHandlers) HandleAddCoverToFavorites(w http.ResponseWriter, r *http.Request) {
	// Получить id юзера

	userID := mux.Vars(r)["user_id"]
	userIDConverted, err := uuid.Parse(userID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Получить id обложки

	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Добавить запись в таблицу

	if err = fh.favoritesService.AddCoverToFavorites(fh.ctx, userIDConverted, coverIDConverted); err != nil {
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (fh *HTTPFavoritesHandlers) HandleGetFavoriteCovers(w http.ResponseWriter, r *http.Request) {
	// Получить id юзера

	userID := mux.Vars(r)["user_id"]
	userIDConverted, err := uuid.Parse(userID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Получить его избранные обложки
	covers, err := fh.favoritesService.GetFavoriteCovers(fh.ctx, userIDConverted)
	if err != nil {
		fmt.Println("ошибка при получении избранных обложек")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("ошибка при получении избранных обложек")
		// Логировать ошибку
	}
}
