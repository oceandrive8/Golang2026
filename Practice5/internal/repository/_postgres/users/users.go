package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"awesomeProject/internal/repository/_postgres"
	"awesomeProject/pkg/modules"
	_ "database/sql"
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
func (r *Repository) CreateUser(u *modules.User) (uuid.UUID, error) {

	u.ID = uuid.New()

	query := `
		INSERT INTO users (id, name, email)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.DB.Exec(query, u.ID, u.Name, u.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return u.ID, nil
}

// ------------------ Read All ------------------
func (r *Repository) GetUsers(limit, offset int, orderBy []string) ([]modules.User, error) {
	var users []modules.User

	orderClause := "id"
	if len(orderBy) > 0 {
		orderClause = strings.Join(orderBy, ", ")
	}

	query := fmt.Sprintf(`
		SELECT id, name, email, gender, birthday
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderClause)

	rows, err := r.db.DB.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u modules.User
		var email sql.NullString
		var birthday sql.NullTime
		var gender sql.NullString

		if err := rows.Scan(&u.ID, &u.Name, &email, &gender, &birthday); err != nil {
			return nil, err
		}

		if email.Valid {
			u.Email = &email.String
		} else {
			u.Email = nil
		}
		if gender.Valid {
			u.Gender = gender.String
		}
		if birthday.Valid {
			u.Birthday = &birthday.Time
		} else {
			u.Birthday = nil
		}

		users = append(users, u)
	}

	return users, nil
}

// ------------------ Read by ID ------------------
func (r *Repository) GetUserByID(id uuid.UUID) (*modules.User, error) {

	var user modules.User

	query := `
		SELECT id, name, email
		FROM users
		WHERE id=$1 AND deleted_at IS NULL
	`

	err := r.db.DB.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("user with ID %s not found: %w", id, err)
	}

	return &user, nil
}

// ------------------ Update ------------------
func (r *Repository) UpdateUser(u *modules.User) error {

	query := `
		UPDATE users
		SET name=$1, email=$2
		WHERE id=$3 AND deleted_at IS NULL
	`

	res, err := r.db.DB.Exec(query, u.Name, u.Email, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no active user found with ID %s", u.ID)
	}

	return nil
}

// ------------------ Delete ------------------
func (r *Repository) DeleteUser(id uuid.UUID) (int64, error) {

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
		return 0, fmt.Errorf("no active user found with ID %s", id)
	}

	return rowsAffected, nil
}

func (r *Repository) GetPaginatedUsers(
	page int,
	pageSize int,
	filters map[string]interface{},
	orderBy []string,
) (modules.PaginatedResponse, error) {

	var users []modules.User
	offset := (page - 1) * pageSize
	args := []interface{}{}
	conditions := []string{}
	i := 1

	// Build WHERE conditions from filters
	for field, value := range filters {
		switch field {
		case "id", "name", "email", "gender", "birthday":
			conditions = append(conditions, fmt.Sprintf("%s = $%d", field, i))
			args = append(args, value)
			i++
		}
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM users"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	var totalCount int
	if err := r.db.DB.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return modules.PaginatedResponse{}, fmt.Errorf("failed to count users: %w", err)
	}

	// Fetch paginated data
	query := "SELECT id, name, email, gender, birthday FROM users"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if len(orderBy) > 0 {
		query += " ORDER BY " + strings.Join(orderBy, ", ")
	} else {
		query += " ORDER BY name ASC"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.DB.Query(query, args...)
	if err != nil {
		return modules.PaginatedResponse{}, fmt.Errorf("failed to fetch paginated users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u modules.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.Birthday); err != nil {
			return modules.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return modules.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (r *Repository) GetCommonFriends(ctx context.Context, userID1, userID2 uuid.UUID) ([]modules.User, error) {
	query := `
    SELECT u.id, u.name, u.email, u.gender, u.birthday
    FROM users u
    JOIN user_friends f1 ON u.id = f1.friend_id
    JOIN user_friends f2 ON u.id = f2.friend_id
    WHERE f1.user_id = $1
      AND f2.user_id = $2
    ORDER BY u.name ASC;
    `

	rows, err := r.db.DB.QueryContext(ctx, query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []modules.User
	for rows.Next() {
		var u modules.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.Birthday); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}

	return friends, nil
}
