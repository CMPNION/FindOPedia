package postgres

import (
	"context"
	"database/sql"
	"findopedia/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 RETURNING id, username, password_hash, created_at`,
		username, passwordHash,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, username, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
