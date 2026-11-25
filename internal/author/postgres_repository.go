package author

import (
	"context"
	"strings"

	"github.com/HemanthGowda/library-management-go/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ExistByName(ctx context.Context, name string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(
		ctx,
		`SELECT EXISTS(
            SELECT 1 FROM authors
            WHERE name = $1
        )`,
		name,
	).Scan(&exists)

	return exists, err
}

func (r *PostgresRepository) Create(ctx context.Context, author Author) (Author, error) {

	author.Id = uuid.NewString()

	query := `
		INSERT INTO authors (
			id,
			name,
			biography,
			birth,
			nationality
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(
		ctx,
		query,
		author.Id,
		author.Name,
		author.Biography,
		author.Birth,
		author.Nationality,
	).Scan(
		&author.Id,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if db.IsUniqueViolation(err) {
		return Author{}, ErrDuplicate
	}

	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func (r *PostgresRepository) ViewAuthor(ctx context.Context, id string) (Author, error) {
	var author Author
	query := `
        SELECT id, name, biography, birth, nationality,
               created_at, updated_at
        FROM authors
        WHERE id = $1
    `
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.Birth,
		&author.Nationality,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func (r *PostgresRepository) SearchAuthor(ctx context.Context, name string, limit, offset int) ([]Author, error) {
	query := `
		SELECT id, name, biography, birth, nationality, created_at, updated_at
		FROM authors
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
		`

	rows, err := r.pool.Query(ctx, query, "%"+strings.TrimSpace(name)+"%", limit, offset)

	if err != nil {
		return nil, err
	}

	var authors []Author

	for rows.Next() {
		var author Author

		err := rows.Scan(
			&author.Id,
			&author.Name,
			&author.Biography,
			&author.Birth,
			&author.Nationality,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *PostgresRepository) CountAuthors(ctx context.Context, name string) (int, error) {

	var count int

	err := r.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM authors
        WHERE name ILIKE $1
    `, "%"+strings.TrimSpace(name)+"%").Scan(&count)

	return count, err
}

func (r *PostgresRepository) HasBooksAuthors(ctx context.Context, id string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
            SELECT 1 FROM books
            WHERE id = $1
        )`, id,
	).Scan(&exists)

	return exists, err
}

func (r *PostgresRepository) DeleteAuthor(ctx context.Context, id string) (Author, error) {
	var author Author

	err := r.pool.QueryRow(ctx,
		`DELETE FROM authors
         WHERE id = $1
         RETURNING id, name, biography, birth, nationality`,
		id,
	).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.Birth,
		&author.Nationality,
	)

	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func (r *PostgresRepository) HasName(ctx context.Context, id string, name string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM authors
			WHERE name = $1 AND id != $2
		)`,
		name, id,
	).Scan(&exists)

	return exists, err
}

func (r *PostgresRepository) UpdateAuthor(ctx context.Context, author Author) (Author, error) {

	err := r.pool.QueryRow(ctx, `
		UPDATE authors
		SET name = $2,
		    biography = $3,
		    birth = $4,
		    nationality = $5,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, biography, birth, nationality, created_at, updated_at
	`,
		author.Id,
		author.Name,
		author.Biography,
		author.Birth,
		author.Nationality,
	).Scan(
		&author.Id,
		&author.Name,
		&author.Biography,
		&author.Birth,
		&author.Nationality,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		return Author{}, err
	}

	return author, nil
}
