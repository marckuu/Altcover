package services

import (
	"book-cover-service/db/repositories"
	"book-cover-service/domains"
	"context"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepository repositories.BookRepository
}

func NewBookService(bookRepository repositories.BookRepository) BookService {
	return BookService{
		bookRepository: bookRepository,
	}
}

func (b *BookService) AddBook(ctx context.Context, book domains.Book) error {
	if err := b.bookRepository.AddBook(ctx, book); err != nil {
		return err
	}
	return nil
}

func (b *BookService) UpdateBook(ctx context.Context, book domains.Book) error {
	if err := b.bookRepository.UpdateBook(ctx, book); err != nil {
		return err
	}
	return nil
}

func (b *BookService) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	if err := b.bookRepository.DeleteBook(ctx, bookID); err != nil {
		return err
	}
	return nil
}
