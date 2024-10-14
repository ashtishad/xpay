package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
)

// CardFilters defines the filters for querying cards, Used in List()
type CardFilters struct {
	UserID   *int64
	WalletID *int64
	Provider *string
	Status   *string
}

type CardRepository interface {
	AddCardToWallet(ctx context.Context, card *Card) (*Card, common.AppError)
	FindBy(ctx context.Context, dbColumnName string, value any) (*Card, common.AppError)
	Update(ctx context.Context, card *Card) common.AppError
	Delete(ctx context.Context, cardID string) common.AppError
	List(ctx context.Context, filters CardFilters) ([]*Card, common.AppError)
}

type cardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{
		db: db,
	}
}

// AddCardToWallet adds a new card to a wallet, using serializable isolation to prevent
// concurrent addition of duplicate cards for the same user and wallet.
// It checks for existing cards before insertion and handles potential conflicts.
func (r *cardRepository) AddCardToWallet(ctx context.Context, card *Card) (*Card, common.AppError) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "AddCardToWallet")

	if appErr := r.checkExistingCard(ctx, tx, card); appErr != nil {
		return nil, appErr
	}

	query := `INSERT INTO cards (uuid, user_id, wallet_id, encrypted_card_number, provider, type, last_four, expiry_date, status)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		card.UUID, card.UserID, card.WalletID, card.EncryptedCardNumber, card.Provider, card.Type,
		card.LastFour, card.ExpiryDate, card.Status).
		Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		slog.Error("failed to add card to wallet", "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return card, nil
}

// FindBy retrieves a card based on a specific database column (id, uuid, wallet_id, user_id),
// using read committed isolation to ensure consistent reads across different lookup methods.
func (r *cardRepository) FindBy(ctx context.Context, dbColumnName string, value any) (*Card, common.AppError) {
	query, err := r.generateFindByQuery(dbColumnName)
	if err != nil {
		return nil, common.NewBadRequestError(common.ErrUnexpectedDatabase).Wrap(err)
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "FindBy")

	var card Card
	err = tx.QueryRowContext(ctx, query, value).Scan(
		&card.ID, &card.UUID, &card.UserID, &card.WalletID, &card.EncryptedCardNumber, &card.Provider, &card.Type,
		&card.LastFour, &card.ExpiryDate, &card.Status, &card.CreatedAt, &card.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NewNotFoundError("card not found")
		}

		slog.Error("failed to get card", "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return &card, nil
}

// Update modifies the mutable fields of a card (expiry date and status), using serializable
// isolation to ensure atomic updates and prevent concurrent modifications.
func (r *cardRepository) Update(ctx context.Context, card *Card) common.AppError {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Update")

	query := `UPDATE cards SET expiry_date = $1, status = $2, updated_at = NOW() WHERE id = $3`

	if _, err = tx.ExecContext(ctx, query, card.ExpiryDate, card.Status, card.ID); err != nil {
		slog.Error("failed to update card", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return nil
}

// Delete performs a soft delete on a card by changing its status to 'deleted', using
// serializable isolation to ensure atomic status updates.
func (r *cardRepository) Delete(ctx context.Context, cardID string) common.AppError {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "Delete")

	query := `UPDATE cards SET status = $1, updated_at = NOW() WHERE uuid = $2`

	result, err := tx.ExecContext(ctx, query, CardStatusDeleted, cardID)
	if err != nil {
		slog.Error("failed to soft delete card", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get rows affected", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if rowsAffected == 0 {
		slog.Warn("zero rows affected")
		return common.NewNotFoundError("card not found")
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return nil
}

// List retrieves multiple cards based on provided filters, using read committed isolation
// to ensure consistent reads while allowing for concurrent transactions.
func (r *cardRepository) List(ctx context.Context, filters CardFilters) ([]*Card, common.AppError) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		slog.Error(common.ErrTXBegin, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer rollBackOnError(tx, "List")

	query, args := r.buildListQuery(filters)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("failed to query cards", "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	defer func(rows *sql.Rows) {
		clsErr := rows.Close()
		if clsErr != nil {
			slog.WarnContext(ctx, "failed to close rows", "err", clsErr)
		}
	}(rows)

	var cards []*Card
	for rows.Next() {
		var card Card
		err := rows.Scan(
			&card.ID, &card.UUID, &card.UserID, &card.WalletID, &card.EncryptedCardNumber, &card.Provider, &card.Type,
			&card.LastFour, &card.ExpiryDate, &card.Status, &card.CreatedAt, &card.UpdatedAt)

		if err != nil {
			slog.Error("failed to scan card", "err", err)
			return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
		}

		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		slog.Error("error iterating over rows", "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err = tx.Commit(); err != nil {
		slog.Error(common.ErrTxCommit, "err", err)
		return nil, common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	return cards, nil
}

// checkExistingCard verifies if a card with the same type and provider already exists for the user,
// preventing duplicate entries and handling reactivation of previously deleted cards.
func (r *cardRepository) checkExistingCard(ctx context.Context, tx *sql.Tx, card *Card) common.AppError {
	var existingCardStatus string
	query := `SELECT status FROM cards WHERE user_id = $1 AND provider = $2 AND type = $3 LIMIT 1`
	err := tx.QueryRowContext(ctx, query, card.UserID, card.Provider, card.Type).Scan(&existingCardStatus)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error("failed to check for existing card", "err", err)
		return common.NewInternalServerError(common.ErrUnexpectedDatabase, err)
	}

	if err == nil {
		if existingCardStatus == CardStatusDeleted {
			return common.NewConflictError("A deleted card of this type and provider exists. Please update its status instead of adding a new one.")
		} else {
			return common.NewConflictError("An active or inactive card of this type and provider already exists.")
		}
	}

	return nil
}

// generateFindByQuery creates the appropriate SQL query based on the specified field name,
// supporting flexible querying for the FindBy method while preventing SQL injection.
func (r *cardRepository) generateFindByQuery(fieldName string) (string, error) {
	baseQuery := `SELECT id, uuid, user_id, wallet_id, encrypted_card_number, provider, type, last_four, expiry_date, status, created_at, updated_at
				  FROM cards WHERE status != 'deleted' AND `

	switch fieldName {
	case common.DBColumnID:
		return baseQuery + "id = $1", nil
	case common.DBColumnUUID:
		return baseQuery + "uuid = $1", nil
	case common.DBColumnUserID:
		return baseQuery + "user_id = $1", nil
	case common.DBColumnWalletID:
		return baseQuery + "wallet_id = $1", nil
	default:
		return "", errors.New("invalid db field name for card lookup")
	}
}

// buildListQuery constructs the SQL query and arguments for listing cards based on the provided CardFilters.
func (r *cardRepository) buildListQuery(filters CardFilters) (string, []any) {
	query := `SELECT id, uuid, user_id, wallet_id, encrypted_card_number, provider, type, last_four, expiry_date, status, created_at, updated_at
              FROM cards
              WHERE 1=1`
	var args []any
	argCount := 1

	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
		argCount++
	}

	if filters.WalletID != nil {
		query += fmt.Sprintf(" AND wallet_id = $%d", argCount)
		args = append(args, *filters.WalletID)
		argCount++
	}

	if filters.Provider != nil {
		query += fmt.Sprintf(" AND provider = $%d", argCount)
		args = append(args, *filters.Provider)
		argCount++
	}

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	} else {
		query += " AND status != 'deleted'"
	}

	query += " ORDER BY created_at DESC"

	return query, args
}
