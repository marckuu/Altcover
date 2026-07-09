package transport

import (
	"book-cover-service/core/domains"
	logs "book-cover-service/core/logger"
	"book-cover-service/core/tools"
	serviceInterfaces "book-cover-service/internal/books/services/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPBookHandlers struct {
	bookService serviceInterfaces.BookService
	ctx         context.Context
	logger      logs.Logger
}

func NewHTTPBookHandlers(bookService serviceInterfaces.BookService, ctx context.Context, logger logs.Logger) HTTPBookHandlers {
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
