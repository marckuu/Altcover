package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type BookService interface {
	AddBook(ctx context.Context, book domains.Book) error
	UpdateBook(ctx context.Context, book domains.Book) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
}
