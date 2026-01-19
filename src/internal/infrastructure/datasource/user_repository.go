package datasource

import (
	"context"
	"database/sql"
	"errors"
	"school21/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
}

type userRepository struct {
	db *DB
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := u.db.Pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash)
			 VALUES ($1, $2, $3)`, user.ID, user.Login, user.PasswordHash)
	return err
}

func (u *userRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := u.db.Pool.QueryRow(ctx,
		`SELECT id, username, password_hash
			 FROM users
			 WHERE username = $1`, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err

}

func NewUserRepository(db *DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
