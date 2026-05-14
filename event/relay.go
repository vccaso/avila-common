package event

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// EventConfig contains basic config for Relay
type EventConfig struct {
	ServiceName string
	MaxRetries  int
	DB          *sql.DB
	JS          nats.JetStreamContext
}

type OutboxRow struct {
	ID            int64
	EventID       string
	CorrelationID sql.NullString
	EventType     string
	AggregateType string
	AggregateID   string
	Payload       string
	CustomerID    string
	Status        string
	RetryCount    int
}

// StartRelayWorker starts a goroutine that continuously polls the outbox and publishes to JetStream.
func StartRelayWorker(ctx context.Context, cfg EventConfig) {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	workerID := fmt.Sprintf("%s:%s:%d:%s", cfg.ServiceName, hostname, pid, uuid.New().String())

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("OutboxRelay worker stopped.")
				return
			case <-ticker.C:
				processOutbox(ctx, cfg, workerID)
			}
		}
	}()
}

func processOutbox(ctx context.Context, cfg EventConfig, workerID string) {
	// 1. Lock rows
	lockQuery := `
		UPDATE outbox_events 
		SET status = 'processing', locked_at = NOW(), locked_by = ? 
		WHERE 
		  status = 'pending' 
		  OR (
			status = 'processing' 
			AND locked_at < DATE_SUB(NOW(), INTERVAL 5 MINUTE)
		  )
		ORDER BY occurred_at ASC 
		LIMIT 100
	`
	res, err := cfg.DB.ExecContext(ctx, lockQuery, workerID)
	if err != nil {
		log.Printf("[%s] Error locking outbox rows: %v", cfg.ServiceName, err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return // Nothing to process
	}

	// 2. Fetch locked rows
	fetchQuery := `
		SELECT id, event_id, correlation_id, event_type, aggregate_type, aggregate_id, payload, customer_id, status, retry_count
		FROM outbox_events
		WHERE status = 'processing' AND locked_by = ?
		ORDER BY occurred_at ASC
	`
	rows, err := cfg.DB.QueryContext(ctx, fetchQuery, workerID)
	if err != nil {
		log.Printf("[%s] Error fetching locked rows: %v", cfg.ServiceName, err)
		return
	}
	defer rows.Close()

	var events []OutboxRow
	for rows.Next() {
		var e OutboxRow
		err := rows.Scan(&e.ID, &e.EventID, &e.CorrelationID, &e.EventType, &e.AggregateType, &e.AggregateID, &e.Payload, &e.CustomerID, &e.Status, &e.RetryCount)
		if err != nil {
			log.Printf("[%s] Error scanning row: %v", cfg.ServiceName, err)
			continue
		}
		events = append(events, e)
	}
	rows.Close()

	// 3. Publish to NATS
	for _, e := range events {
		subject := fmt.Sprintf("eazychat.%s", e.EventType)
		
		correlationID := ""
		if e.CorrelationID.Valid {
			correlationID = e.CorrelationID.String
		}

		// Construct standard JSON payload envelope
		finalPayload := fmt.Sprintf(`{
			"event_id": "%s",
			"correlation_id": "%s",
			"customer_id": "%s",
			"source_service": "%s",
			"aggregate_type": "%s",
			"aggregate_id": "%s",
			"event_type": "%s",
			"data": %s
		}`, e.EventID, correlationID, e.CustomerID, cfg.ServiceName, e.AggregateType, e.AggregateID, e.EventType, e.Payload)

		msg := &nats.Msg{
			Subject: subject,
			Data:    []byte(finalPayload),
		}

		_, err := cfg.JS.PublishMsg(msg)
		if err != nil {
			log.Printf("[%s] Failed to publish event %s: %v", cfg.ServiceName, e.EventID, err)
			markFailed(ctx, cfg.DB, e.ID, e.RetryCount, cfg.MaxRetries, err.Error())
		} else {
			markPublished(ctx, cfg.DB, e.ID)
		}
	}
}

func markPublished(ctx context.Context, db *sql.DB, id int64) {
	_, _ = db.ExecContext(ctx, "UPDATE outbox_events SET status = 'published', published_at = NOW() WHERE id = ?", id)
}

func markFailed(ctx context.Context, db *sql.DB, id int64, currentRetry, maxRetries int, errMsg string) {
	newRetry := currentRetry + 1
	var status string

	if newRetry >= maxRetries {
		status = "dead_letter"
	} else {
		status = "pending"
	}

	_, _ = db.ExecContext(ctx, "UPDATE outbox_events SET status = ?, retry_count = ?, last_error = ?, locked_by = NULL, locked_at = NULL WHERE id = ?", 
		status, newRetry, errMsg, id)
}
