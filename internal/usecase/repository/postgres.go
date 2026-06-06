package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/project/library/internal/entity"
	"go.uber.org/zap"
)

const pgFKViolation = "23503"

func isFKViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgFKViolation
}

type postgresRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresRepository(db *pgxpool.Pool, logger *zap.Logger) *postgresRepository {
	return &postgresRepository{
		db:     db,
		logger: logger,
	}

}

func (p *postgresRepository) AddBook(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.logger.Info("failed to begin transaction", zap.String("book id", book.ID))
		return nil, err
	}
	defer tx.Rollback(ctx)
	const queryBook = `INSERT INTO book (id, name) VALUES($1, $2) RETURNING created_at, updated_at`
	err = tx.QueryRow(ctx, queryBook, book.ID, book.Name).Scan(&book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return nil, err
	}

	const queryAuthorBook = `INSERT INTO author_book (author_id, book_id) VALUES($1, $2)`
	for _, authorID := range book.AuthorIDs {
		_, err = tx.Exec(ctx, queryAuthorBook, authorID, book.ID)
		if err != nil {
			if isFKViolation(err) {
				return nil, entity.ErrAuthorNotFound
			}
			return nil, err
		}
	}
	const queryFetchAuthors = `
		SELECT COALESCE(array_agg(author_id ORDER BY author_id), '{}')
		FROM author_book
		WHERE book_id = $1
	`
	if err := tx.QueryRow(ctx, queryFetchAuthors, book.ID).Scan(&book.AuthorIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return book, nil
}

func (s *postgresRepository) GetBook(ctx context.Context, ID string) (*entity.Book, error) {

	const query = `
    SELECT book.id, book.name, book.created_at, book.updated_at,
        COALESCE(array_agg(author_book.author_id ORDER BY author_book.author_id) FILTER (WHERE author_book.author_id IS NOT NULL), '{}') AS author_ids
    FROM book
    LEFT JOIN author_book ON book.id = author_book.book_id
    WHERE book.id = $1
    GROUP BY book.id
`
	var book entity.Book
	err := s.db.QueryRow(ctx, query, ID).Scan(
		&book.ID, &book.Name, &book.CreatedAt, &book.UpdatedAt, &book.AuthorIDs,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (p *postgresRepository) UpdateBook(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.logger.Info("failed to begin transaction", zap.String("book id", book.ID))
		return nil, err
	}
	defer tx.Rollback(ctx)

	const queryName = `UPDATE book SET name = $1 WHERE id = $2`
	tag, err := tx.Exec(ctx, queryName, book.Name, book.ID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, entity.ErrBookNotFound
	}

	const queryDelete = `delete from author_book where book_id = $1`
	_, err = tx.Exec(ctx, queryDelete, book.ID)
	if err != nil {
		return nil, err
	}

	const queryAuthorBook = `INSERT INTO author_book (author_id, book_id) VALUES($1, $2)`
	for _, authorID := range book.AuthorIDs {
		_, err = tx.Exec(ctx, queryAuthorBook, authorID, book.ID)
		if err != nil {
			if isFKViolation(err) {
				return nil, entity.ErrAuthorNotFound
			}
			return nil, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (p *postgresRepository) RegisterAuthor(ctx context.Context, author *entity.Author) (*entity.Author, error) {
	const queryInsertAuthor = `insert into author(id, name, created_at, updated_at) values($1, $2, $3, $4)`
	_, err := p.db.Exec(ctx, queryInsertAuthor, author.ID, author.Name, author.CreatedAt, author.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (p *postgresRepository) GetAuthor(ctx context.Context, id string) (*entity.Author, error) {
	const queryGetAuthor = `SELECT id, name, created_at, updated_at FROM author WHERE id = $1`
	res := &entity.Author{}
	err := p.db.QueryRow(ctx, queryGetAuthor, id).
		Scan(&res.ID, &res.Name, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrAuthorNotFound
		}
		return nil, err
	}
	return res, nil
}

func (p *postgresRepository) UpdateAuthor(ctx context.Context, author *entity.Author) (*entity.Author, error) {
	const queryUpdateAuthor = `UPDATE author SET name = $1 WHERE id = $2 RETURNING created_at, updated_at`
	err := p.db.QueryRow(ctx, queryUpdateAuthor, author.Name, author.ID).
		Scan(&author.CreatedAt, &author.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrAuthorNotFound
		}
		return nil, err
	}
	return author, nil
}

func (p *postgresRepository) GetAuthorBooks(ctx context.Context, authorID string) ([]*entity.Book, error) {
	const query = `
        SELECT b.id, b.name, b.created_at, b.updated_at,
               COALESCE(array_agg(ab2.author_id) FILTER (WHERE ab2.author_id IS NOT NULL), '{}')
        FROM book b
        JOIN author_book ab ON ab.book_id = b.id
        LEFT JOIN author_book ab2 ON ab2.book_id = b.id
        WHERE ab.author_id = $1
        GROUP BY b.id`

	rows, err := p.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.Book
	for rows.Next() {
		b := &entity.Book{}
		if err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.AuthorIDs); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}
