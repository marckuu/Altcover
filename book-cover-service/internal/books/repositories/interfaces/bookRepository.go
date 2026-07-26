package interfaces

import (
	"book-cover-service/core/domains"
	"context"

	"github.com/google/uuid"
)

type BookRepository interface {
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
	AddBook(ctx context.Context, book domains.Book) (domains.Book, error)
	UpdateBook(ctx context.Context, book domains.Book) (domains.Book, error)
	GetBookByTitle(ctx context.Context, title string) (domains.Book, error)
}
