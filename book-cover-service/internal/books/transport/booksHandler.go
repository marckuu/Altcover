package transport

import (
	"book-cover-service/core/domains"
	"book-cover-service/core/tools"
	"book-cover-service/internal/books/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type HTTPBookHandlers struct {
	bookService services.BookService
	ctx         context.Context
	logger      *zap.Logger
}

func NewHTTPBookHandlers(bookService services.BookService, ctx context.Context, logger *zap.Logger) HTTPBookHandlers {
	return HTTPBookHandlers{
		bookService: bookService,
		ctx:         ctx,
		logger:      logger,
	}
}

func (b *HTTPBookHandlers) HandleAddBook(w http.ResponseWriter, r *http.Request) {
	var book domains.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := b.bookService.AddBook(b.ctx, book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при сохранении книги: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (b *HTTPBookHandlers) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var book domains.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if err := b.bookService.UpdateBook(b.ctx, book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при сохранении книги: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (b *HTTPBookHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	bookIDRaw := mux.Vars(r)["book_id"]
	bookID, err := uuid.Parse(bookIDRaw)
	if err != nil {
		b.logger.Error(fmt.Errorf("не удалось преобразовать полученный id книги в uuid: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = b.bookService.DeleteBook(b.ctx, bookID); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при удалении книги: %w", err).Error())
		tools.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
