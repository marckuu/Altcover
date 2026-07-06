package transport

import (
	"auth-service/core/domains"
	"auth-service/core/errors"
	logs "auth-service/core/logger"
	"auth-service/core/middleware"
	"auth-service/internal/designerProfiles/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPDesignerProfileHandlers struct {
	designerProfileService services.DesignerProfileService
	ctx                    context.Context
	logger                 logs.Logger
}

func NewHTTPDesignerProfileHandlers(designerProfileService services.DesignerProfileService,
	ctx context.Context,
	logger logs.Logger) HTTPDesignerProfileHandlers {
	return HTTPDesignerProfileHandlers{
		designerProfileService: designerProfileService,
		ctx:                    ctx,
		logger:                 logger,
	}
}

func (d *HTTPDesignerProfileHandlers) HandleCreateMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать переданный идентификатор: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	_, err = d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err == nil {
		d.logger.Error(fmt.Errorf("профиль дизайнера уже существует: %w", err).Error())
		errors.SendErrorResponse(w, fmt.Errorf("профиль текущего дизайнера уже существует"), http.StatusInternalServerError)
		return
	}

	var designerProfile domains.DesignerProfile

	if err = json.NewDecoder(r.Body).Decode(&designerProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать тело запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	designerProfile.UserID = userID

	if err = d.designerProfileService.CreateDesignerProfile(d.ctx, designerProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось создать профиль дизайнера: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (d *HTTPDesignerProfileHandlers) HandleUpdateMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	var newDesignerProfile domains.DesignerProfile

	if err := json.NewDecoder(r.Body).Decode(&newDesignerProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать тело запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать переданный идентификатор: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	oldDesignerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		d.logger.Error(fmt.Errorf("ошибка получения профиля дизайнера: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	newDesignerProfile.ID = oldDesignerProfile.ID

	if err = d.designerProfileService.UpdateDesignerProfile(d.ctx, newDesignerProfile); err != nil {
		d.logger.Error(fmt.Errorf("ошибка при сохранении обложки: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (d *HTTPDesignerProfileHandlers) HandleGetMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать переданный идентификатор: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	designerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		d.logger.Error(fmt.Errorf("ошибка получения профиля дизайнера: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(designerProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось отправить ответ с профилем дизайнера: %w", err).Error())
	}
}

func (d *HTTPDesignerProfileHandlers) HandleDeleteMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось получить ID пользователя из токена: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	designerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось получить профиль дизайнера: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = d.designerProfileService.DeleteDesignerProfile(d.ctx, designerProfile.ID); err != nil {
		d.logger.Error(fmt.Errorf("не удалось удалить профиль: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
