package domain

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
)

// UserRepository defines the interface for user data operations.
// It abstracts the underlying data interactions, allowing for flexible implementations.
type UserRepository interface {
	FindIDFromUUID(ctx context.Context, uuid string) (int64, common.AppError)
	Create(ctx context.Context, user *User) (*User, common.AppError)
	FindBy(ctx context.Context, dbColumnName string, value any) (*User, common.AppError)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// FindIDFromUUID retrieves a user's ID by UUID. Useful for referencing users in other tables.
// Returns NotFoundError if user doesn't exist, or InternalServerError on database errors.
func (r *userRepository) FindIDFromUUID(ctx context.Context, uuid string) (int64, common.AppError) {
	query := `SELECT id FROM users WHERE uuid = $1`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, common.NewNotFoundError("user not found by uuid")
		}

		slog.Error("failed to get user ID", "uuid", uuid, "err", err)
		return 0, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return userID, nil
}

// Create inserts a new user using a serializable transaction. Checks for email uniqueness.
// Returns the created user with ID on success, or AppError (409 for email conflict, 500 for other errors).
func (r *userRepository) Create(ctx context.Context, u *User) (*User, common.AppError) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Create User")

	if appError := r.checkUserExistsByEmail(ctx, tx, u.Email); appError != nil {
		return nil, appError
	}

	createdID, appError := r.insertUser(ctx, tx, u)
	if appError != nil {
		return nil, appError
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	u.ID = createdID
	return u, nil
}

// checkUserExistsByEmail verifies email uniqueness within the Create transaction.
// Returns ConflictError if email exists, or InternalServerError on database errors.
func (r *userRepository) checkUserExistsByEmail(ctx context.Context, tx *sql.Tx, email string) common.AppError {
	var exists bool

	existsQuery := `SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)`
	if err := tx.QueryRowContext(ctx, existsQuery, email).Scan(&exists); err != nil {
		slog.Error("failed to check user exist by email", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if exists {
		return common.NewConflictError("user with this email already exists")
	}

	return nil
}

// insertUser performs the actual user insertion within the Create transaction.
// Returns the new user's ID or InternalServerError on failure.
func (r *userRepository) insertUser(ctx context.Context, tx *sql.Tx, u *User) (int64, common.AppError) {
	queryCreateUser := `INSERT INTO users (uuid, full_name, email, password_hash, status, role, created_at, updated_at)
                        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
                        RETURNING id`

	var createdID int64
	err := tx.QueryRowContext(ctx, queryCreateUser,
		u.UUID, u.FullName, u.Email, u.PasswordHash, u.Status, u.Role, u.CreatedAt, u.UpdatedAt).Scan(&createdID)

	if err != nil {
		slog.Error("failed to create user", "err", err)
		return 0, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return createdID, nil
}

// FindBy retrieves a user by id, uuid, or email.
// Returns the user or appropriate AppError (NotFoundError or InternalServerError).
func (r *userRepository) FindBy(ctx context.Context, dbColumnName string, value any) (*User, common.AppError) {
	query, err := generateFindByQuery(dbColumnName)
	if err != nil || query == "" {
		return nil, common.NewBadRequestError(common.ErrUnexpectedDatabase)
	}

	var user User
	err = r.db.QueryRowContext(ctx, query, value).Scan(
		&user.ID, &user.UUID, &user.FullName, &user.Email, &user.PasswordHash,
		&user.Status, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewNotFoundError("user not found")
		}

		slog.Error("failed to get user", "field", dbColumnName, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return &user, nil
}

// generateFindByQuery creates SQL query for FindBy method, supporting id, uuid, and email fields.
// Returns the query string or an error for invalid db field.
func generateFindByQuery(fieldName string) (string, error) {
	baseQuery := `SELECT id, uuid, full_name, email, password_hash, status, role, created_at, updated_at
                  FROM users WHERE `

	var condition string
	switch fieldName {
	case common.DBColumnID:
		condition = "id = $1"
	case common.DBColumnUUID:
		condition = "uuid = $1"
	case common.DBColumnEmail:
		condition = "email = $1"
	default:
		return "", errors.New("invalid db field name")
	}

	return baseQuery + condition, nil
}
