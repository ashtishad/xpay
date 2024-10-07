package domain

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
)

// UserRepository defines the interface for user data operations.
// It abstracts the underlying data interactions, allowing for flexible implementations.
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, common.AppError)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database using a serializable transaction.
// Checks if email exists (returns 409 Conflict if true), Inserts new user if email is unique
// Returns: Created user with ID on success, AppError on failure.
// Error codes: 409 for email conflict, 500 for other errors.
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

// checkUserExistsByEmail verifies if a user with the given email already exists.
// It's part of the Create transaction to prevent duplicate entries.
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

// insertUser performs the actual insertion of a new user into the database.
// It's called within the Create transaction after checking for existing emails.
// Returns the ID of the newly created user or an InternalServerError if the insertion fails.
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
