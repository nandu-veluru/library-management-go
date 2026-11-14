package author

import (
    "context"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/HemanthGowda/library-management-go/internal/db"
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

    err := r.pool.QueryRow(
        ctx,
        `INSERT INTO authors (id, name, biography, birth, nationality)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id`,
        author.Id,
        author.Name,
        author.Biography,
        author.Birth,
        author.Nationality,
    ).Scan(&author.Id)

    if db.IsUniqueViolation(err) {
        return Author{}, ErrDuplicate
    }

    if err != nil {
        return Author{}, err
    }

    return author, nil
}