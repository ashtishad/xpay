package domain

import (
	"database/sql"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
)

func rollBackOnError(tx *sql.Tx, methodName string) {
	if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
		slog.Error(common.ErrTXRollback, "err", rbErr, "method", methodName)
	}
}
