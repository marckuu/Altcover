package handlers

import (
	"book-cover-service/handlers/tools"
	"book-cover-service/middleware"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPFavoritesHandlers struct {
	favoritesService services.FavoritesService
	ctx              context.Context
}

func NewHTTPFavoritesHandlers(favoriteService services.FavoritesService, ctx context.Context) HTTPFavoritesHandlers {
	return HTTPFavoritesHandlers{
		favoritesService: favoriteService,
		ctx:              ctx,
	}
}

func (f *HTTPFavoritesHandlers) HandleAddCoverToFavorites(w http.ResponseWriter, r *http.Request) {
	// Получить id юзера

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Получить id обложки

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Добавить запись в таблицу

	if err = f.favoritesService.AddCoverToFavorites(f.ctx, userID, coverID); err != nil {
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (f *HTTPFavoritesHandlers) HandleGetMyFavoriteCovers(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		fmt.Println("ошибка получения offset из query параметра")
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		fmt.Println("ошибка получения limit из query параметра")
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Получить id юзера

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Получить его избранные обложки
	covers, err := f.favoritesService.GetFavoriteCovers(f.ctx, userID, int(offset), int(limit))
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
