package transport

import (
	"book-cover-service/core/domains"
	"book-cover-service/core/errors"
	logs "book-cover-service/core/logger"
	serviceInterfaces "book-cover-service/internal/books/services/interfaces"
	"book-cover-service/internal/books/transport/dto"
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

// @Summary Add book
// @Description Создать книгу
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Param request body dto.CreateBookRequest true "Book data"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books [post]
func (b *HTTPBookHandlers) HandleAddBook(w http.ResponseWriter, r *http.Request) {
	var createBookRequest dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&createBookRequest); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	book := domains.Book{
		Title:       createBookRequest.Title,
		Description: createBookRequest.Description,
	}
	if err := b.bookService.AddBook(b.ctx, book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при сохранении книги: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при записи ответа с созданной книгой: %w", err).Error())
	}
}

// @Summary Update book
// @Description Обновить книгу
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Param request body dto.UpdateBookRequest true "Book data"
// @Succes 200
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /books [patch]
func (b *HTTPBookHandlers) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBookRequest dto.UpdateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&updateBookRequest); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при чтении тела запроса: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	book := domains.Book{
		ID:          updateBookRequest.ID,
		Title:       updateBookRequest.Title,
		Description: updateBookRequest.Description,
	}
	if err := b.bookService.UpdateBook(b.ctx, book); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при сохранении книги: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

// @Summary Delete book
// @Description Удалить книгу
// @Security ApiKeyAuth
// @Tags books
// @Param book_id path string true "book id"
// @Succes 200
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/{book_id} [delete]
func (b *HTTPBookHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	bookIDRaw := mux.Vars(r)["book_id"]
	bookID, err := uuid.Parse(bookIDRaw)
	if err != nil {
		b.logger.Error(fmt.Errorf("не удалось преобразовать полученный id книги в uuid: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = b.bookService.DeleteBook(b.ctx, bookID); err != nil {
		b.logger.Error(fmt.Errorf("ошибка при удалении книги: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

// @Summary Get book by title
// @Description Получить книгу по названию
// @Security ApiKeyAuth
// @Tags books
// @Param title path string true "book title"
// @Succes 200
// @Failure 500 {object} errors.ErrorResponse
// @Router /books/{title} [get]
func (b *HTTPBookHandlers) HandleGetBookByTitle(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	book, err := b.bookService.GetBookByTitle(b.ctx, title)
	if err != nil {
		b.logger.Error(fmt.Errorf("ошибка при удалении книги: %w", err).Error())
		errors.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(book); err != nil {
		b.logger.Error(fmt.Errorf("не удалось записать ответ с полученной книгой: %w", err).Error())
	}
}
