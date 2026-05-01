package handlers

import (
	"book-cover-service/domains"
	"book-cover-service/handlers/tools"
	"book-cover-service/services"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPCoverHandlers struct {
	coverService services.CoverService
	ctx          context.Context
}

func (ch *HTTPCoverHandlers) HandleGetCoversByDesignerID(w http.ResponseWriter, r *http.Request, offset int, limit int) {
	designerID := mux.Vars(r)["designer_id"]
	designerIDConverted, err := uuid.Parse(designerID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	covers, err := ch.coverService.GetCoversByDesignerID(ch.ctx, offset, limit, designerIDConverted)
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

func (ch *HTTPCoverHandlers) HandleUpdateCover(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	// Нужно проверить, что пользователь действительно автор этой обложки
	designerID := mux.Vars(r)["designer_id"]
	designerIDConverted, err := uuid.Parse(designerID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := ch.coverService.GetCoverByID(ctx, coverIDConverted)
	if err != nil {
		fmt.Println("Ошибка при проверке авторства")
		return
	}

	if cover.DesignerID != designerIDConverted {
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
	if err = ch.coverService.UpdateCover(ctx, newCover); err != nil {
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

func (ch *HTTPCoverHandlers) HandleGetCoverByID(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	coverID := mux.Vars(r)["cover_id"]
	coverIDConverted, err := uuid.Parse(coverID)
	if err != nil {
		fmt.Println("ошибка преобразования строки в uuid")
		tools.SendErrorResponse(w, errors.New("не получилось сконвертировать id"), http.StatusInternalServerError)
		return
	}

	cover, err := ch.coverService.GetCoverByID(ctx, coverIDConverted)
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

func (ch *HTTPCoverHandlers) HandleAddCover(w http.ResponseWriter, r *http.Request) {
	// Считать данные обложки из тела запроса
	var cover domains.Cover

	if err := json.NewDecoder(r.Body).Decode(&cover); err != nil {
		fmt.Printf("Ошибка при чтении тела запроса: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Провалидировать обложку

	// Записать обложку через репозиторий
	if err := ch.coverService.AddCover(ch.ctx, cover); err != nil {
		fmt.Printf("Ошибка при сохранении обложки: %s", err)
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Вернуть в ответ добавленную обложку
}

func (ch *HTTPCoverHandlers) HandleGetFeedCovers() {
	// Получить список обложк с наибольшим количеством поставленных лайков за последние три дня

	// Нужно отслеживать количество лайков за последние три дня
}

//func (ch *HTTPCoverHandlers) HandleSetLikeToCover(w http.ResponseWriter, r *http.Request) {
//	designerID := mux.Vars(r)["cover_id"]
//
//	// Нужно как-то понимать был ли поставлен лайк
//
//}
