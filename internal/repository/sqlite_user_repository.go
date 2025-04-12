package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/watanabenaoki/go-clean-arch/internal/domain/model"
)

type sqliteUserRepository struct {
	db *sql.DB
}


func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepository{
		db: db,
	}
}

func (r *sqliteUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	user := &model.User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		log.Printf("Error getting user by id: %v", err)
		return nil, model.ErrInternalServerError
	}

	return user, nil
}

func (r *sqliteUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
	}

	return user, nil
}


// Create implements UserRepository.
func (r *sqliteUserRepository) Create(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}

// Delete implements UserRepository.
func (r *sqliteUserRepository) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// List implements UserRepository.
func (r *sqliteUserRepository) List(ctx context.Context, limit int, offset int) ([]*model.User, error) {
	panic("unimplemented")
}

// Update implements UserRepository.
func (r *sqliteUserRepository) Update(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}
