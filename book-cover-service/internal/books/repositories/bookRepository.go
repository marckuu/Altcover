package repositories

import (
	"book-cover-service/core/db"
	"book-cover-service/core/domains"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errBookNotFound = errors.New("книга не найдена")

type BookRepository struct {
	connection db.Database
}

func NewBookRepository(conn db.Database) BookRepository {
	return BookRepository{
		connection: conn,
	}
}

func (b *BookRepository) AddBook(ctx context.Context, book domains.Book) (domains.Book, error) {
	query := `
	INSERT INTO book (title, description)
	VALUES ($1, $2)
	RETURNING id, title, description;
`
	resultRow := b.connection.QueryRow(ctx, query, book.Title, book.Description)

	var savedBook domains.Book

	if err := resultRow.Scan(&savedBook.ID, &savedBook.Title, &savedBook.Description); err != nil {
		return domains.Book{}, fmt.Errorf("book repo / get by id: %w", err)
	}

	return savedBook, nil
}

func (b *BookRepository) GetBookByID(ctx context.Context, bookID int64) (domains.Book, error) {
	query := `
	SELECT id, title, description
	FROM book
	WHERE id = $1;
`
	resultRow := b.connection.QueryRow(ctx, query, bookID)

	var book domains.Book

	if err := resultRow.Scan(&book.ID, &book.Title, &book.Description); err != nil {
		return domains.Book{}, fmt.Errorf("book repo / get by id: %w", err)
	}

	return book, nil
}

func (b *BookRepository) GetBookByTitle(ctx context.Context, title string) (domains.Book, error) {
	query := `
	SELECT id, title, description
	FROM book
	WHERE title = $1;
`
	resultRow := b.connection.QueryRow(ctx, query, title)

	var book domains.Book

	if err := resultRow.Scan(&book.ID, &book.Title, &book.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Book{}, errBookNotFound
		}
		return domains.Book{}, fmt.Errorf("book repo / get by title: %w", err)
	}

	return book, nil
}

func (b *BookRepository) GetBooks(ctx context.Context, offset int, limit int) ([]domains.Book, error) {
	query := `
	SELECT id, title, description
	FROM book
	OFFSET $1
	LIMIT $2;
`
	resultRows, err := b.connection.Query(ctx, query, offset, limit)
	if err != nil {
		return []domains.Book{}, fmt.Errorf("book repo / get all / query: %w", err)
	}

	defer resultRows.Close()

	var books []domains.Book

	for resultRows.Next() {
		var book domains.Book

		if err = resultRows.Scan(&book.ID, &book.Title, &book.Description); err != nil {
			return []domains.Book{}, fmt.Errorf("book repo / get all / parsing: %w", err)
		}

		books = append(books, book)
	}

	if resultRows.Err() != nil {
		return []domains.Book{}, fmt.Errorf("book repo / get all / query result: %w", err)
	}

	return books, nil
}

func (b *BookRepository) UpdateBook(ctx context.Context, book domains.Book) (domains.Book, error) {
	query := `
	UPDATE book
	SET title = $1, description = $2
	WHERE id = $3
	RETURNING id, title, description;
`
	resultRow := b.connection.QueryRow(ctx, query, book.Title, book.Description, book.ID)

	var savedBook domains.Book

	if err := resultRow.Scan(&savedBook.ID, &savedBook.Title, &savedBook.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Book{}, errBookNotFound
		}
		return domains.Book{}, fmt.Errorf("book repo / update book: %w", err)
	}

	return savedBook, nil
}

func (b *BookRepository) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	query := `
	DELETE FROM book
	WHERE id = $1;
`
	tag, err := b.connection.Exec(ctx, query, bookID)
	if err != nil {
		return fmt.Errorf("book repo / delete: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errBookNotFound
	}

	return nil
}
