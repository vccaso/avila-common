package event

import (
	"context"
	"database/sql"
)

// HandleEventFunc is the business logic function signature
type HandleEventFunc func(ctx context.Context, tx *sql.Tx, msgData []byte) error

// ProcessIdempotentEvent handles the duplicate check, transaction management, and business logic execution.
// It returns true if the message should be Ack'd (success or duplicate), and false + error if it should be Nak'd.
func ProcessIdempotentEvent(ctx context.Context, db *sql.DB, eventID string, msgData []byte, handler HandleEventFunc) (bool, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// 1. Insert Ignore for idempotency check
	res, err := tx.ExecContext(ctx, "INSERT IGNORE INTO processed_events (event_id, processed_at) VALUES (?, NOW())", eventID)
	if err != nil {
		return false, err // DB Error, might need to retry, so don't Ack yet
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// Duplicate! We rollback the business transaction (already deferred), but we MUST return success to Ack the message
		return true, nil
	}

	// 2. Execute Business Logic
	err = handler(ctx, tx, msgData)
	if err != nil {
		return false, err // Business logic failed, return error to Nak the message
	}

	// 3. Commit
	if err = tx.Commit(); err != nil {
		return false, err // Commit failed, Nak
	}

	return true, nil // Success, Ack
}
