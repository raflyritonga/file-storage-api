package repository

import (
	"context"
	"database/sql"

	"github.com/raflyritonga/file-storage-api/internal/model"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
        INSERT INTO users (id, email, username, password, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Username, user.Password, user.Role)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)
	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}	
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := r.db.QueryRowContext(ctx, query, username)
	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}	
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, username, password, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}	
	return &user, nil
}