package transport

import (
	"auth-service/core/domains"
	"auth-service/core/errors"
	logs "auth-service/core/logger"
	"auth-service/core/middleware"
	"auth-service/internal/designerProfiles/services/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPDesignerProfileHandlers struct {
	designerProfileService interfaces.DesignerProfileService
	ctx                    context.Context
	logger                 logs.Logger
}

func NewHTTPDesignerProfileHandlers(designerProfileService interfaces.DesignerProfileService,
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

	var designerProfile domains.DesignerProfile

	if err = json.NewDecoder(r.Body).Decode(&designerProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать тело запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = d.designerProfileService.CreateDesignerProfileToUser(d.ctx, userID, designerProfile); err != nil {
		d.logger.Error(err.Error())
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

	if err = d.designerProfileService.UpdateDesignerProfileToUser(d.ctx, userID, newDesignerProfile); err != nil {
		d.logger.Error(err.Error())
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

	if err = d.designerProfileService.DeleteDesignerProfileToUser(d.ctx, userID); err != nil {
		d.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
