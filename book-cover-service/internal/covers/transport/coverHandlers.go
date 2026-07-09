package transport

import (
	"book-cover-service/core/domains"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	tools2 "book-cover-service/core/tools"
	serviceInterfaces "book-cover-service/internal/covers/services/interfaces"
	designerProfileSnapshotServiceInterfaces "book-cover-service/internal/snapshots/services/interfaces"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var mostLikedInterval = 3

type HTTPCoverHandlers struct {
	coverService                   serviceInterfaces.CoverService
	designerProfileSnapshotService designerProfileSnapshotServiceInterfaces.DesignerProfileSnapshotService
	ctx                            context.Context
	logger                         logs.Logger
}

func NewCoverHandlers(coverService serviceInterfaces.CoverService,
	designerProfileSnapshotService designerProfileSnapshotServiceInterfaces.DesignerProfileSnapshotService,
	ctx context.Context,
	logger logs.Logger) HTTPCoverHandlers {
	return HTTPCoverHandlers{
		coverService:                   coverService,
		designerProfileSnapshotService: designerProfileSnapshotService,
		ctx:                            ctx,
		logger:                         logger,
	}
}

func (c *HTTPCoverHandlers) HandleGetMyCoversAsDesigner(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools2.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, offset, limit, userID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении списка обложек дизайнера: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

func (c *HTTPCoverHandlers) HandleGetCoversByDesignerUserID(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools2.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	userIDRaw := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, offset, limit, userID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении списка обложек дизайнера: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

func (c *HTTPCoverHandlers) HandleUpdateCover(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		tools2.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := c.coverService.GetCoverByID(c.ctx, coverID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при проверке авторства: %w", err).Error())
		return
	}

	if cover.UserID != userID {
		c.logger.Error(fmt.Errorf("дизайнер не является автором данной обложки: %w", err).Error())
		tools2.SendErrorResponse(w, errors.New("дизайнер не является автором обложки"), http.StatusConflict)
		return
	}

	var newCover domains.Cover

	if err = json.NewDecoder(r.Body).Decode(&newCover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении новой обложки из запроса: %w", err).Error())
		tools2.SendErrorResponse(w, errors.New("ошибка при получении новой обложки из запроса"), http.StatusBadRequest)
		return
	}

	newCover.ID = coverID
	newCover.UserID = userID
	newCover.BookID = cover.BookID

	if err = c.coverService.UpdateCover(c.ctx, newCover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при обновлении обложки: %w", err).Error())
		tools2.SendErrorResponse(w, errors.New("ошибка при обновлении обложки"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(newCover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

func (c *HTTPCoverHandlers) HandleGetCoverByID(w http.ResponseWriter, r *http.Request) {
	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		tools2.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := c.coverService.GetCoverByID(c.ctx, coverIDConverted)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении обложки: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(cover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученной обложкой: %w", err).Error())
	}
}

func (c *HTTPCoverHandlers) HandleAddCover(w http.ResponseWriter, r *http.Request) {
	var cover domains.Cover

	if err := json.NewDecoder(r.Body).Decode(&cover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	cover.UserID = userID

	profileSnapshot, err := c.designerProfileSnapshotService.GetDesignerProfileSnapshotByUserID(c.ctx, userID)
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить профиль пользователя: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	cover.DesignerNickname = profileSnapshot.Nickname
	cover.DesignerAvatarKey = profileSnapshot.AvatarKey

	// Провалидировать обложку

	if err = c.coverService.AddCover(c.ctx, cover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при сохранении обложки: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Вернуть в ответ добавленную обложку
}

func (c *HTTPCoverHandlers) HandleGetFeedCovers(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools2.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	covers, err := c.coverService.GetMostLikedCovers(c.ctx, mostLikedInterval, offset, limit)
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить список популярных обложек: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("не удалось отправить ответ со списком популярных обложек: %w", err).Error())
	}
}

func (c *HTTPCoverHandlers) HandleGetCoversByBook(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools2.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		tools2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	bookIDRaw := mux.Vars(r)["book_id"]
	bookID, err := uuid.Parse(bookIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	covers, err := c.coverService.GetCoversByBook(c.ctx, bookID, offset, limit)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка получения списка обложек по книге: %w", err).Error())
		tools2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("не удалось отправить ответ со списком обложек по книге: %w", err).Error())
	}
}
