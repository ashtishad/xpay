package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
)

// WalletRepository defines the interface for wallet data operations.
type WalletRepository interface {
	FindIDFromUUID(ctx context.Context, walletUUID string) (int64, common.AppError)
	Create(ctx context.Context, wallet *Wallet) (*Wallet, common.AppError)
	UpdateStatus(ctx context.Context, walletUUID string, status string) common.AppError
	FindBy(ctx context.Context, dbColumnName string, value any) (*Wallet, common.AppError)
	GetBalance(ctx context.Context, walletUUID string) (int64, common.AppError)
}

type walletRepository struct {
	db *sql.DB
}

// NewWalletRepository creates a new instance of WalletRepository.
func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

// FindIDFromUUID retrieves a wallet's ID by UUID. Useful for referencing users in other tables.
// Returns NotFoundError if user doesn't exist, or InternalServerError on database errors.
func (r *walletRepository) FindIDFromUUID(ctx context.Context, uuid string) (int64, common.AppError) {
	query := `SELECT id FROM wallets WHERE uuid = $1`

	var walletID int64
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(&walletID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, common.NewNotFoundError("wallet not found by uuid")
		}

		slog.Error("failed to get wallet ID", "uuid", uuid, "err", err)
		return 0, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return walletID, nil
}

// Create inserts a new wallet into the database.
// It uses a serializable transaction to prevent race conditions and ensure data consistency.
// Before insertion, it checks for existing wallets for the same user and currency to prevent duplicates.
func (r *walletRepository) Create(ctx context.Context, wallet *Wallet) (*Wallet, common.AppError) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}
	defer rollBackOnError(tx, "Create Wallet")

	if appErr := r.checkExistingWallet(ctx, tx, wallet.UserID, wallet.Currency); appErr != nil {
		return nil, appErr
	}

	createdWallet, appErr := r.insertWallet(ctx, tx, wallet)
	if appErr != nil {
		return nil, appErr
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return createdWallet, nil
}

// UpdateStatus changes the status of a wallet.
// It uses a serializable transaction to ensure atomic updates and prevent conflicts.
// The method returns a NotFoundError if the wallet doesn't exist.
func (r *walletRepository) UpdateStatus(ctx context.Context, walletUUID string, status string) common.AppError {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Update Wallet Status")

	query := `UPDATE wallets SET status = $1 WHERE uuid = $2`

	result, err := tx.ExecContext(ctx, query, status, walletUUID)
	if err != nil {
		slog.Error("failed to update wallet status", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if rowsAffected == 0 {
		slog.Warn("row affected is zero")
		return common.NewNotFoundError("wallet not found")
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return nil
}

// FindBy retrieves a wallet based on the specified column and value.
// It uses read committed isolation to ensure consistent reads while allowing for better concurrency than serializable.
// The method supports lookups by ID, UUID, and UserID (for active wallets only).
func (r *walletRepository) FindBy(ctx context.Context, dbColumnName string, value any) (*Wallet, common.AppError) {
	query, err := r.generateFindByQuery(dbColumnName)
	if err != nil || query == "" {
		slog.Error("failed to generate FindBy query", "column", dbColumnName)
		return nil, common.NewBadRequestError(common.ErrUnexpectedDatabase).Wrap(err)
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Find Wallet By...")

	var wallet Wallet
	err = tx.QueryRowContext(ctx, query, value).Scan(
		&wallet.ID, &wallet.UUID, &wallet.UserID, &wallet.BalanceInCents, &wallet.Currency,
		&wallet.Status, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewNotFoundError("wallet not found")
		}
		slog.Error("failed to get wallet", "field", dbColumnName, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return &wallet, nil
}

// GetBalance retrieves the current balance of a wallet given its UUID.
// It uses a READ COMMITTED transaction to ensure consistent reads.
func (r *walletRepository) GetBalance(ctx context.Context, walletUUID string) (int64, common.AppError) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return 0, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Get Wallet Balance")

	query := `SELECT balance FROM wallets WHERE uuid = $1 AND status = 'active'`

	var balance int64
	err = tx.QueryRowContext(ctx, query, walletUUID).Scan(&balance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, common.NewNotFoundError("Wallet not found or wallet is not active")
		}

		slog.Error("failed to get wallet balance", "err", err, "uuid", walletUUID)
		return 0, common.NewInternalServerError("Failed to get wallet balance", err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return 0, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return balance, nil
}

// checkExistingWallet is a helper method for the Create operation.
// It checks if a wallet already exists for the given user and currency.
// This prevents duplicate wallets and provides a clear error message if a wallet already exists.
func (r *walletRepository) checkExistingWallet(ctx context.Context, tx *sql.Tx, userID int64, currency string) common.AppError {
	query := `SELECT uuid, status FROM wallets WHERE user_id = $1 AND currency = $2`
	var existingWalletUUID, existingWalletStatus string
	err := tx.QueryRowContext(ctx, query, userID, currency).Scan(&existingWalletUUID, &existingWalletStatus)

	if err == nil {
		return common.NewConflictError(fmt.Sprintf("user already has a wallet (UUID: %s, Status: %s) for this currency.Please update the wallet status if needed", existingWalletUUID, existingWalletStatus))
	}

	if err != sql.ErrNoRows {
		slog.Error("failed to check existing wallet", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return nil
}

// insertWallet is a helper method for the Create operation.
// It performs the actual insertion of the wallet into the database.
// It handles unique constraint violations and returns appropriate errors.
func (r *walletRepository) insertWallet(ctx context.Context, tx *sql.Tx, wallet *Wallet) (*Wallet, common.AppError) {
	query := `INSERT INTO wallets (uuid, user_id, balance, currency, status)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, query,
		wallet.UUID, wallet.UserID, wallet.BalanceInCents, wallet.Currency, wallet.Status).
		Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		slog.Error("failed to create wallet", "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return wallet, nil
}

// generateFindByQuery is a helper method for the FindBy operation.
// It generates the appropriate SQL query based on the column name provided.
// This method centralizes query generation and helps prevent SQL injection.
func (r *walletRepository) generateFindByQuery(fieldName string) (string, error) {
	baseQuery := `SELECT id, uuid, user_id, balance, currency, status, created_at, updated_at
				  FROM wallets WHERE `

	switch fieldName {
	case common.DBColumnID:
		return baseQuery + "id = $1", nil
	case common.DBColumnUUID:
		return baseQuery + "uuid = $1", nil
	case common.DBColumnUserID:
		return baseQuery + "user_id = $1 AND status = 'active'", nil
	default:
		return "", errors.New("invalid db field name for wallet lookup")
	}
}
