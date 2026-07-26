package transport

import (
	"book-cover-service/core/domains"
	errors2 "book-cover-service/core/errors"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/middleware"
	tools "book-cover-service/core/tools"
	serviceInterfaces "book-cover-service/internal/covers/services/interfaces"
	"book-cover-service/internal/covers/transport/dto"
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

// @Summary Get my covers
// @Description Получить обложки текущего дизайнера
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param offset query string true "Pagination offset"
// @Param limit query string true "Pagination limit"
// @Succes 200 {array} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /designers/me/covers [get]
func (c *HTTPCoverHandlers) HandleGetMyCoversAsDesigner(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, offset, limit, userID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении списка обложек дизайнера: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

// @Summary Get covers by user ID
// @Description Получить обложки по id пользователя
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param offset query string true "Pagination offset"
// @Param limit query string true "Pagination limit"
// @Param user_id path string true "User id"
// @Succes 200 {array} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /designers/{user_id}/covers [get]
func (c *HTTPCoverHandlers) HandleGetCoversByDesignerUserID(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	userIDRaw := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	covers, err := c.coverService.GetCoversByUserID(c.ctx, offset, limit, userID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении списка обложек дизайнера: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

// @Summary Update cover
// @Description Обновить обложку
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param cover_id path string true "Cover id"
// @Param request body dto.UpdateCoverRequest true "New cover"
// @Succes 200 {object} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /covers/{cover_id} [patch]
func (c *HTTPCoverHandlers) HandleUpdateCover(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coverIDRaw := mux.Vars(r)["cover_id"]
	coverID, err := uuid.Parse(coverIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors2.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	var updateCoverRequest dto.UpdateCoverRequest

	if err = json.NewDecoder(r.Body).Decode(&updateCoverRequest); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	cover := domains.Cover{
		Title:       updateCoverRequest.Title,
		Description: updateCoverRequest.Description,
		ImagesKeys:  updateCoverRequest.ImagesKeys,
		Status:      updateCoverRequest.Status,
	}
	newCover, err := c.coverService.UpdateCover(c.ctx, coverID, userID, cover)
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(newCover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученными обложками: %w", err).Error())
	}
}

// @Summary Get cover by id
// @Description Получить обложку по id
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param cover_id path string true "Cover id"
// @Succes 200 {object} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /covers/{cover_id} [get]
func (c *HTTPCoverHandlers) HandleGetCoverByID(w http.ResponseWriter, r *http.Request) {
	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors2.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := c.coverService.GetCoverByID(c.ctx, coverIDConverted)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка при получении обложки: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(cover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с полученной обложкой: %w", err).Error())
	}
}

// @Summary Add cover
// @Description Добавить обложку
// @Security ApiKeyAuth
// @Tags covers
// @Accept json
// @Produce json
// @Param request body dto.CreateCoverRequest true "Cover id"
// @Succes 200 {object} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /covers [post]
func (c *HTTPCoverHandlers) HandleAddCover(w http.ResponseWriter, r *http.Request) {
	var createCoverRequest dto.CreateCoverRequest

	if err := json.NewDecoder(r.Body).Decode(&createCoverRequest); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	cover := domains.Cover{
		Title:       createCoverRequest.Title,
		Description: createCoverRequest.Description,
		ImagesKeys:  createCoverRequest.ImagesKeys,
		Status:      createCoverRequest.Status,
		UserID:      userID,
		BookID:      createCoverRequest.BookID,
	}
	savedCover, err := c.coverService.AddCover(c.ctx, cover)
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(savedCover); err != nil {
		c.logger.Error(fmt.Errorf("ошибка при записи ответа с созданной обложкой: %w", err).Error())
	}
}

// @Summary Get feed covers
// @Description Получить обложки для ленты рекомендаций
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param offset query string true "Pagination offset"
// @Param limit query string true "Pagination limit"
// @Succes 200 {array} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /feeds/covers [get]
func (c *HTTPCoverHandlers) HandleGetFeedCovers(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	covers, err := c.coverService.GetMostLikedCovers(c.ctx, mostLikedInterval, offset, limit)
	if err != nil {
		c.logger.Error(fmt.Errorf("не удалось получить список популярных обложек: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if covers == nil {
		covers, err = c.coverService.GetCovers(c.ctx, offset, limit)
		if err != nil {
			c.logger.Error(fmt.Errorf("не удалось получить список всех обложек: %w", err).Error())
			errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("не удалось отправить ответ со списком популярных обложек: %w", err).Error())
	}
}

// @Summary Get covers by book
// @Description Получить обложки книги
// @Security ApiKeyAuth
// @Tags covers
// @Produce json
// @Param book_id path string true "Book id"
// @Param offset query string true "Pagination offset"
// @Param limit query string true "Pagination limit"
// @Succes 200 {array} domains.Cover
// @Failure 400 {object} errors2.ErrorResponse
// @Failure 500 {object} errors2.ErrorResponse
// @Router /books/{book_id}/covers [get]
func (c *HTTPCoverHandlers) HandleGetCoversByBook(w http.ResponseWriter, r *http.Request) {
	offset, limit, err := tools.GetOffsetAndLimitFromQuery(r.URL.Query().Get("offset"), r.URL.Query().Get("limit"))
	if err != nil {
		c.logger.Error(err.Error())
		errors2.SendErrorResponse(w, err, http.StatusBadRequest)
	}

	bookIDRaw := mux.Vars(r)["book_id"]
	bookID, err := uuid.Parse(bookIDRaw)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка преобразования строки в uuid: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	covers, err := c.coverService.GetCoversByBook(c.ctx, bookID, offset, limit)
	if err != nil {
		c.logger.Error(fmt.Errorf("ошибка получения списка обложек по книге: %w", err).Error())
		errors2.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(covers); err != nil {
		c.logger.Error(fmt.Errorf("не удалось отправить ответ со списком обложек по книге: %w", err).Error())
	}
}
