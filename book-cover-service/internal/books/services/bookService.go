package services

import (
	"book-cover-service/core/domains"
	repositoryInterfaces "book-cover-service/internal/books/repositories/interfaces"
	"context"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepository repositoryInterfaces.BookRepository
}

func NewBookService(bookRepository repositoryInterfaces.BookRepository) *BookService {
	return &BookService{
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

func (b *BookService) GetBookByTitle(ctx context.Context, title string) (domains.Book, error) {
	book, err := b.bookRepository.GetBookByTitle(ctx, title)
	if err != nil {
		return domains.Book{}, err
	}
	return book, nil
}
