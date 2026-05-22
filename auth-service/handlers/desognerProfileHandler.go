package handlers

import (
	"auth-service/domains"
	"auth-service/handlers/tools"
	"auth-service/middleware"
	"auth-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPDesignerProfileHandlers struct {
	designerProfileService services.DesignerProfileService
	ctx                    context.Context
}

func NewHTTPDesignerProfileHandlers() HTTPDesignerProfileHandlers {
	return HTTPDesignerProfileHandlers{
		designerProfileService: services.DesignerProfileService{},
		ctx:                    nil,
	}
}

func (d *HTTPDesignerProfileHandlers) HandleCreateDesignerProfile(w http.ResponseWriter, r *http.Request) {
	// Считать новый профиль и тела
	var designerProfile domains.DesignerProfile

	if err := json.NewDecoder(r.Body).Decode(&designerProfile); err != nil {
		fmt.Printf("Не удалось распознать тело запроса: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Провалидировать профиль

	if err := d.designerProfileService.CreateDesignerProfile(d.ctx, designerProfile); err != nil {
		fmt.Printf("Не удалось создать профиль дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (d *HTTPDesignerProfileHandlers) HandleUpdateCurrentDesignerProfile(w http.ResponseWriter, r *http.Request) {
	// Считать новый профиль и тела
	var newDesignerProfile domains.DesignerProfile

	if err := json.NewDecoder(r.Body).Decode(&newDesignerProfile); err != nil {
		fmt.Printf("Не удалось распознать тело запроса: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось распознать переданный идентификатор: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	oldDesignerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		fmt.Printf("Ошибка получения профиля дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	newDesignerProfile.ID = oldDesignerProfile.ID

	// Провалидировать профиль

	// Обновить профиль
	if err = d.designerProfileService.UpdateDesignerProfile(d.ctx, newDesignerProfile); err != nil {
		fmt.Printf("Ошибка при сохранении обложки: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

}

func (d *HTTPDesignerProfileHandlers) HandleGetCurrentDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось распознать переданный идентификатор: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Найти профиль и вернуть его
	designerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		fmt.Printf("Ошибка получения профиля дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Записать профиль в ответ
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(designerProfile); err != nil {
		fmt.Printf("Не удалось отправить ответ с профилем дизайнера: %s", err)
		// Логировать ошибку
	}
}

func (d *HTTPDesignerProfileHandlers) HandleDeleteCurrentDesignerProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		fmt.Printf("Не удалось получить ID пользователя из токена: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	designerProfile, err := d.designerProfileService.GetProfileByUserID(d.ctx, userID)
	if err != nil {
		fmt.Printf("Не удалось получить профиль дизайнера: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = d.designerProfileService.DeleteDesignerProfile(d.ctx, designerProfile.ID); err != nil {
		fmt.Printf("Не удалось удалить профиль: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
