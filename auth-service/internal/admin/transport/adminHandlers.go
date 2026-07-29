package transport

import (
	"auth-service/core/enums"
	"auth-service/core/errors"
	logs "auth-service/core/logger"
	"auth-service/internal/admin/services/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPAdminHandlers struct {
	adminService interfaces.AdminService
	ctx          context.Context
	logger       logs.Logger
}

func NewHTTPAdminHandlers(
	adminService interfaces.AdminService,
	ctx context.Context,
	logger logs.Logger) HTTPAdminHandlers {
	return HTTPAdminHandlers{
		adminService: adminService,
		ctx:          ctx,
		logger:       logger,
	}
}

// @Summary Change role
// @Description Изменяет роль пользователя
// @Tags admin
// @Security ApiKeyAuth
// @Accept json
// @Param user_id path string true "User id"
// @Success 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /role/{user_id} [patch]
func (a *HTTPAdminHandlers) HandleChangeRole(w http.ResponseWriter, r *http.Request) {
	userIDRaw := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		a.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return

	}

	var userRole enums.Role
	if err = json.NewDecoder(r.Body).Decode(&userRole); err != nil {
		a.logger.Error(fmt.Errorf("не удалось получить роль из запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = a.adminService.ChangeRole(a.ctx, userID, userRole); err != nil {
		a.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

// @Summary Get role
// @Description Возвращает роль пользователя
// @Tags admin
// @Security ApiKeyAuth
// @Accept json
// @Param user_id path string true "User id"
// @Success 200
// @Failure 400 {object} errors.ErrorResponse
// @Router /role/{user_id} [get]
func (a *HTTPAdminHandlers) HandleGetRole(w http.ResponseWriter, r *http.Request) {
	userIDRaw := mux.Vars(r)["user_id"]
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		a.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	role, err := a.adminService.GetRole(a.ctx, userID)
	if err != nil {
		a.logger.Error(err.Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(role); err != nil {
		a.logger.Error(err.Error())
	}
}
