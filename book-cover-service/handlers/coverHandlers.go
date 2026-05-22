package handlers

import (
	"book-cover-service/domains"
	"book-cover-service/handlers/tools"
	"book-cover-service/middleware"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var mostLikedInterval = 3

type HTTPCoverHandlers struct {
	coverService services.CoverService
	ctx          context.Context
}

func NewCoverHandlers(coverService services.CoverService, ctx context.Context) HTTPCoverHandlers {
	return HTTPCoverHandlers{
		coverService: coverService,
		ctx:          ctx,
	}
}

func (c *HTTPCoverHandlers) HandleGetMyCoversAsDesigner(w http.ResponseWriter, r *http.Request) {
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

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, int(offset), int(limit), userID)
	if err != nil {
		fmt.Printf("Ошибка при получении списка обложек дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("Ошибка при записи ответа с полученными обложками")
		// Логировать оишбку
	}
}

func (c *HTTPCoverHandlers) HandleGetCoversByDesignerUserID(w http.ResponseWriter, r *http.Request) {
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

	userIDRaw := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, int(offset), int(limit), userID)
	if err != nil {
		fmt.Printf("Ошибка при получении списка обложек дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("Ошибка при записи ответа с полученными обложками")
		// Логировать оишбку
	}
}

func (c *HTTPCoverHandlers) HandleUpdateCover(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := c.coverService.GetCoverByID(c.ctx, coverID)
	if err != nil {
		fmt.Println("Ошибка при проверке авторства")
		return
	}

	if cover.UserID != userID {
		fmt.Println("дизайнер не является автором данной обложки")
		tools.SendErrorResponse(w, errors.New("дизайнер не является автором обложки"), http.StatusConflict)
		return
	}

	// Считать обложку
	var newCover domains.Cover

	if err = json.NewDecoder(r.Body).Decode(&newCover); err != nil {
		fmt.Println("ошибка при получении новой обложки из запроса")
		tools.SendErrorResponse(w, errors.New("ошибка при получении новой обложки из запроса"), http.StatusBadRequest)
		return
	}

	// Обновить обложку
	if err = c.coverService.UpdateCover(c.ctx, newCover); err != nil {
		fmt.Println("ошибка при обновлении обложки")
		tools.SendErrorResponse(w, errors.New("ошибка при обновлении обложки"), http.StatusInternalServerError)
		return
	}

	// Вернуть обновленную обложку
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(newCover); err != nil {
		fmt.Println("Ошибка при записи ответа с полученными обложками")
		// Логировать оишбку
	}
}

func (c *HTTPCoverHandlers) HandleGetCoverByID(w http.ResponseWriter, r *http.Request) {
	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := c.coverService.GetCoverByID(c.ctx, coverIDConverted)
	if err != nil {
		fmt.Printf("Ошибка при получении обложки: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(cover); err != nil {
		fmt.Println("Ошибка при записи ответа с полученной обложкой")
		// Логировать оишбку
	}
}

func (c *HTTPCoverHandlers) HandleAddCover(w http.ResponseWriter, r *http.Request) {
	// Считать данные обложки из тела запроса
	var cover domains.Cover

	if err := json.NewDecoder(r.Body).Decode(&cover); err != nil {
		fmt.Printf("Ошибка при чтении тела запроса: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Сделать запрос на получение профиля этого пользователя и заполнить их
	_ = userID

	// Провалидировать обложку

	// Записать обложку через репозиторий
	if err := c.coverService.AddCover(c.ctx, cover); err != nil {
		fmt.Printf("Ошибка при сохранении обложки: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Вернуть в ответ добавленную обложку
}

func (c *HTTPCoverHandlers) HandleGetFeedCovers(w http.ResponseWriter, r *http.Request) {
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

	covers, err := c.coverService.GetMostLikedCovers(c.ctx, mostLikedInterval, int(offset), int(limit))
	if err != nil {
		fmt.Println("не удалось получить список популярных обложек")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("не удалось отправить ответ со списком популярных обложек")
		// Логировать ошибку
	}
}

func (c *HTTPCoverHandlers) HandleGetCoversByBook(w http.ResponseWriter, r *http.Request) {
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

	bookIDRaw := mux.Vars(r)["book_id"]
	bookID, err := uuid.Parse(bookIDRaw)

	covers, err := c.coverService.GetCoversByBook(c.ctx, bookID, int(offset), int(limit))
	if err != nil {
		fmt.Println("ошибка получения списка обложек по книге")
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		fmt.Println("не удалось отправить ответ со списком обложек по книге")
		// Логировать ошибку
	}
}
