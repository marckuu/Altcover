package handlers

import (
	"auth-service/domains"
	"auth-service/handlers/tools"
	"auth-service/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPDesignerProfileHandlers struct {
	designerProfileService services.DesignerProfileService
	ctx                    context.Context
}

func (dh *HTTPDesignerProfileHandlers) UpdateDesignerProfile(w http.ResponseWriter, r *http.Request) {
	// Считать новый профиль и тела
	var designerProfile domains.DesignerProfile

	if err := json.NewDecoder(r.Body).Decode(&designerProfile); err != nil {
		fmt.Printf("Не удалось распознать тело запроса: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Провалидировать профиль

	// Обновить профиль
	if err := dh.designerProfileService.UpdateDesignerProfile(dh.ctx, designerProfile); err != nil {
		fmt.Printf("Ошибка при сохранении обложки: %s", err)
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

}

func (dh *HTTPDesignerProfileHandlers) GetDesignerProfileByUserID(w http.ResponseWriter, r *http.Request) {
	// Получить userID
	userID := mux.Vars(r)["user_id"]

	userIDConverted, err := uuid.Parse(userID)
	if err != nil {
		fmt.Printf("Не удалось распознать переданный идентификатор: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Найти профиль и вернуть его
	designerProfile, err := dh.designerProfileService.GetProfileByUserID(dh.ctx, userIDConverted)
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
