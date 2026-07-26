package transport

import (
	"auth-service/core/domains"
	"auth-service/core/errors"
	logs "auth-service/core/logger"
	"auth-service/core/middleware"
	"auth-service/internal/designerProfiles/services/interfaces"
	"auth-service/internal/designerProfiles/transport/dto"
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

// @Summary Create my designer profile
// @Description Создать профиль дизайнера для текущего пользователя
// @Tags designer profiles
// @Security ApiKeyAuth
// Accept json
// @Param request body dto.CreteDesignerProfileRequest true "Designer profile"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /designer_profiles/me [post]
func (d *HTTPDesignerProfileHandlers) HandleCreateMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать переданный идентификатор: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	var createDesignerProfileRequest dto.CreteDesignerProfileRequest

	if err = json.NewDecoder(r.Body).Decode(&createDesignerProfileRequest); err != nil {
		d.logger.Error(fmt.Errorf("не удалось распознать тело запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	designerProfile := domains.DesignerProfile{
		Nickname:  createDesignerProfileRequest.Nickname,
		AvatarKey: createDesignerProfileRequest.AvatarKey,
		UserID:    userID,
	}
	savedProfile, err := d.designerProfileService.CreateDesignerProfileToUser(d.ctx, userID, designerProfile)
	if err != nil {
		d.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(savedProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось отправить ответ с профилем дизайнера: %w", err).Error())
	}
}

// @Summary Update designer profile
// @Description Обновляет профиль дизайнера текущего пользователя
// @Tags designer profiles
// @Security ApiKeyAuth
// @Accept json
// @Param request body dto.UpdateDesignerProfileRequest true "Designer profile"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /designer_profiles/me [patch]
func (d *HTTPDesignerProfileHandlers) HandleUpdateMyDesignerProfile(w http.ResponseWriter, r *http.Request) {
	var updateDesignerProfileRequest dto.UpdateDesignerProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&updateDesignerProfileRequest); err != nil {
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

	designerProfile := domains.DesignerProfile{
		Nickname:  updateDesignerProfileRequest.Nickname,
		AvatarKey: updateDesignerProfileRequest.Nickname,
	}
	savedProfile, err := d.designerProfileService.UpdateDesignerProfileToUser(d.ctx, userID, designerProfile)
	if err != nil {
		d.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(savedProfile); err != nil {
		d.logger.Error(fmt.Errorf("не удалось отправить ответ с профилем дизайнера: %w", err).Error())
	}
}

// @Summary Get my designer profile
// @Description Получить профиль дизайнера текущего пользователя
// @Tags designer profiles
// @Security ApiKeyAuth
// @Produce json
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /designer_profiles/me [get]
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

// @Summary Delete my designer profile
// @Description Удаляет профиль дизайнера текущего пользователя
// @Tags designer profiles
// @Security ApiKeyAuth
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /designer_profiles/me [delete]
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
