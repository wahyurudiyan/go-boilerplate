package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	userEnt "github.com/wahyurudiyan/go-boilerplate/core/entities/user"
)

// Ensure userRepositoryImpl implements IUserRepository interface
var _ IUserRepository = (*userRepositoryImpl)(nil)

// userRepositoryImpl implements the IUserRepository interface
type userRepositoryImpl struct {
	db *sqlx.DB
}

// NewIUserRepository creates a new instance of IUserRepository
func NewUserSQLRepository(db *sqlx.DB) IUserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// SaveUser inserts a single user into the database
func (r *userRepositoryImpl) SaveUser(ctx context.Context, user userEnt.User) error {
	query := `
		INSERT INTO users (
			role, email, unique_id, fullname, username, password, created_at, updated_at
		) VALUES (
			:role, :email, :unique_id, :fullname, :username, :password, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// SaveUsers inserts multiple users into the database
func (r *userRepositoryImpl) SaveUsers(ctx context.Context, users []userEnt.User) error {
	if len(users) == 0 {
		return nil
	}

	query := `
		INSERT INTO users (
			role, email, unique_id, fullname, username, password, created_at, updated_at
		) VALUES (
			:role, :email, :unique_id, :fullname, :username, :password, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, users)
	if err != nil {
		return fmt.Errorf("failed to save users: %w", err)
	}
	return nil
}

// UpdateUser updates an existing user in the database
func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user userEnt.User) error {
	query := `
		UPDATE users SET
			role = :role,
			email = :email,
			unique_id = :unique_id,
			fullname = :fullname,
			username = :username,
			password = :password,
			updated_at = :updated_at
		WHERE id = :id
	`
	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", user.Id)
	}

	return nil
}

// DeleteUserById deletes a user by ID
func (r *userRepositoryImpl) DeleteUserById(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// DeleteUserByEmail deletes a user by email
func (r *userRepositoryImpl) DeleteUserByEmail(ctx context.Context, email string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE email = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to delete user by email: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with email %s not found", email)
	}

	return nil
}

// DeleteUserByUniqueId deletes a user by unique ID
func (r *userRepositoryImpl) DeleteUserByUniqueId(ctx context.Context, uniqueId string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE unique_id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, uniqueId)
	if err != nil {
		return fmt.Errorf("failed to delete user by unique id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with unique id %s not found", uniqueId)
	}

	return nil
}

// RetrieveAllUser retrieves all users with pagination
func (r *userRepositoryImpl) RetrieveAllUser(ctx context.Context, offset, limit int) ([]userEnt.User, error) {
	var users []userEnt.User
	query := `
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all users: %w", err)
	}
	return users, nil
}

// RetrieveUserById retrieves a user by ID
func (r *userRepositoryImpl) RetrieveUserById(ctx context.Context, id int64) (userEnt.User, error) {
	var user userEnt.User
	query := `
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return userEnt.User{}, fmt.Errorf("user with id %d not found", id)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by id: %w", err)
	}
	return user, nil
}

// RetrieveUserByIds retrieves users by IDs
func (r *userRepositoryImpl) RetrieveUserByIds(ctx context.Context, ids []int64) ([]userEnt.User, error) {
	if len(ids) == 0 {
		return []userEnt.User{}, nil
	}

	// Convert ids to interface{} slice for query args
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// Create placeholders like $1, $2, $3...
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ", "))

	var users []userEnt.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by ids: %w", err)
	}
	return users, nil
}

// RetrieveUserByEmail retrieves a user by email
func (r *userRepositoryImpl) RetrieveUserByEmail(ctx context.Context, email string) (userEnt.User, error) {
	var user userEnt.User
	query := `
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return userEnt.User{}, fmt.Errorf("user with email %s not found", email)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by email: %w", err)
	}
	return user, nil
}

// RetrieveUserByEmails retrieves users by emails
func (r *userRepositoryImpl) RetrieveUserByEmails(ctx context.Context, emails []string) ([]userEnt.User, error) {
	if len(emails) == 0 {
		return []userEnt.User{}, nil
	}

	// Convert emails to interface{} slice for query args
	args := make([]interface{}, len(emails))
	for i, email := range emails {
		args[i] = email
	}

	// Create placeholders like $1, $2, $3...
	placeholders := make([]string, len(emails))
	for i := range emails {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE email IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ", "))

	var users []userEnt.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by emails: %w", err)
	}
	return users, nil
}

// RetrieveUserByUniqueId retrieves a user by unique ID
func (r *userRepositoryImpl) RetrieveUserByUniqueId(ctx context.Context, uniqueId string) (userEnt.User, error) {
	var user userEnt.User
	query := `
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE unique_id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &user, query, uniqueId)
	if err != nil {
		if err == sql.ErrNoRows {
			return userEnt.User{}, fmt.Errorf("user with unique id %s not found", uniqueId)
		}
		return userEnt.User{}, fmt.Errorf("failed to retrieve user by unique id: %w", err)
	}
	return user, nil
}

// RetrieveUserByUniqueIds retrieves users by unique IDs
func (r *userRepositoryImpl) RetrieveUserByUniqueIds(ctx context.Context, uniqueIds []string) ([]userEnt.User, error) {
	if len(uniqueIds) == 0 {
		return []userEnt.User{}, nil
	}

	// Convert uniqueIds to interface{} slice for query args
	args := make([]interface{}, len(uniqueIds))
	for i, uniqueId := range uniqueIds {
		args[i] = uniqueId
	}

	// Create placeholders like $1, $2, $3...
	placeholders := make([]string, len(uniqueIds))
	for i := range uniqueIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT id, role, email, unique_id, fullname, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE unique_id IN (%s) AND deleted_at IS NULL
	`, strings.Join(placeholders, ", "))

	var users []userEnt.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users by unique ids: %w", err)
	}
	return users, nil
}
