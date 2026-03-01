package users

import (
	"fmt"
	"log"
	"time"

	"awesomeProject/internal/repository/_postgres"
	"awesomeProject/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: 5 * time.Second,
	}
}

// ------------------ Create ------------------
func (r *Repository) CreateUser(u *modules.User) (int, error) {
	var id int
	query := `
		INSERT INTO users (name, email, age, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	err := r.db.DB.Get(&id, query, u.Name, u.Email, u.Age, u.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

// ------------------ Read All ------------------
func (r *Repository) GetUsers(limit, offset int) ([]modules.User, error) {
	var users []modules.User

	query := `
        SELECT id, name, email, age, created_at
        FROM users
        WHERE deleted_at IS NULL
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	if err := r.db.DB.Select(&users, query, limit, offset); err != nil {
		log.Println("DB error in GetUsers:", err)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

// ------------------ Read by ID ------------------
func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User

	query := `
		SELECT id, name, email, age, created_at
		FROM users
		WHERE id=$1 AND deleted_at IS NULL
	`

	err := r.db.DB.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("user with ID %d not found: %w", id, err)
	}

	return &user, nil
}

// ------------------ Update ------------------
func (r *Repository) UpdateUser(u *modules.User) error {
	query := `
		UPDATE users
		SET name=$1, email=$2, age=$3
		WHERE id=$4 AND deleted_at IS NULL
	`
	res, err := r.db.DB.Exec(query, u.Name, u.Email, u.Age, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no active user found with ID %d", u.ID)
	}

	return nil
}

// ------------------ Delete ------------------
func (r *Repository) DeleteUser(id int) (int64, error) {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id=$1 AND deleted_at IS NULL
	`

	res, err := r.db.DB.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to soft delete user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("cannot get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("no active user found with ID %d", id)
	}

	return rowsAffected, nil
}
