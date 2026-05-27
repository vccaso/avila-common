package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
			log.Printf("[%s] Successfully published event %s to subject %s", cfg.ServiceName, e.EventID, subject)
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

// OutboxEvent represents an event in the outbox_events table.
type OutboxEvent struct {
	ID            int64      `json:"id"`
	EventID       string     `json:"event_id"`
	CorrelationID *string    `json:"correlation_id"`
	EventType     string     `json:"event_type"`
	AggregateType string     `json:"aggregate_type"`
	AggregateID   string     `json:"aggregate_id"`
	Payload       string     `json:"payload"` // JSON string
	CustomerID    string     `json:"customer_id"`
	Status        string     `json:"status"`
	RetryCount    int        `json:"retry_count"`
	LockedAt      *time.Time `json:"locked_at"`
	LockedBy      *string    `json:"locked_by"`
	LastError     *string    `json:"last_error"`
	OccurredAt    time.Time  `json:"occurred_at"`
	PublishedAt   *time.Time `json:"published_at"`
}

// ProcessedEvent represents a row in processed_events.
type ProcessedEvent struct {
	EventID     string    `json:"event_id"`
	ProcessedAt time.Time `json:"processed_at"`
}

// PaginationMetadata represents pagination information.
type PaginationMetadata struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// PaginatedResponse wraps items with pagination metadata.
type PaginatedResponse[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// QueryOutboxEvents queries the outbox_events table with filtering and pagination.
func QueryOutboxEvents(ctx context.Context, db *sql.DB, customerID, status, eventType, eventID, correlationID string, page, limit int) ([]OutboxEvent, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := "SELECT id, event_id, correlation_id, event_type, aggregate_type, aggregate_id, payload, customer_id, status, retry_count, locked_at, locked_by, last_error, occurred_at, published_at FROM outbox_events WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM outbox_events WHERE 1=1"
	var args []any
	var countArgs []any

	if customerID != "" {
		query += " AND customer_id = ?"
		countQuery += " AND customer_id = ?"
		args = append(args, customerID)
		countArgs = append(countArgs, customerID)
	}
	if status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
		countArgs = append(countArgs, status)
	}
	if eventType != "" {
		query += " AND event_type LIKE ?"
		countQuery += " AND event_type LIKE ?"
		args = append(args, "%"+eventType+"%")
		countArgs = append(countArgs, "%"+eventType+"%")
	}
	if eventID != "" {
		query += " AND event_id = ?"
		countQuery += " AND event_id = ?"
		args = append(args, eventID)
		countArgs = append(countArgs, eventID)
	}
	if correlationID != "" {
		query += " AND correlation_id = ?"
		countQuery += " AND correlation_id = ?"
		args = append(args, correlationID)
		countArgs = append(countArgs, correlationID)
	}

	var total int
	err := db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY occurred_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		var corr sql.NullString
		var lockedAt sql.NullTime
		var lockedBy sql.NullString
		var lastErr sql.NullString
		var pubAt sql.NullTime
		var payloadBytes []byte

		err := rows.Scan(
			&e.ID, &e.EventID, &corr, &e.EventType, &e.AggregateType,
			&e.AggregateID, &payloadBytes, &e.CustomerID, &e.Status,
			&e.RetryCount, &lockedAt, &lockedBy, &lastErr, &e.OccurredAt, &pubAt,
		)
		if err != nil {
			return nil, 0, err
		}

		e.Payload = string(payloadBytes)

		if corr.Valid {
			e.CorrelationID = &corr.String
		}
		if lockedAt.Valid {
			e.LockedAt = &lockedAt.Time
		}
		if lockedBy.Valid {
			e.LockedBy = &lockedBy.String
		}
		if lastErr.Valid {
			e.LastError = &lastErr.String
		}
		if pubAt.Valid {
			e.PublishedAt = &pubAt.Time
		}

		events = append(events, e)
	}

	return events, total, nil
}

// QueryProcessedEvents queries the processed_events table with filtering and pagination.
func QueryProcessedEvents(ctx context.Context, db *sql.DB, eventID string, page, limit int) ([]ProcessedEvent, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := "SELECT event_id, processed_at FROM processed_events WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM processed_events WHERE 1=1"
	var args []any
	var countArgs []any

	if eventID != "" {
		query += " AND event_id = ?"
		countQuery += " AND event_id = ?"
		args = append(args, eventID)
		countArgs = append(countArgs, eventID)
	}

	var total int
	err := db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query += " ORDER BY processed_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []ProcessedEvent
	for rows.Next() {
		var e ProcessedEvent
		err := rows.Scan(&e.EventID, &e.ProcessedAt)
		if err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}

	return events, total, nil
}

// RetryEvent resets status of a specific outbox event to 'pending'.
func RetryEvent(ctx context.Context, db *sql.DB, id int64) (*OutboxEvent, error) {
	updateQuery := `
		UPDATE outbox_events
		SET status = 'pending', retry_count = 0, last_error = NULL, locked_by = NULL, locked_at = NULL
		WHERE id = ?
	`
	_, err := db.ExecContext(ctx, updateQuery, id)
	if err != nil {
		return nil, err
	}

	query := "SELECT id, event_id, correlation_id, event_type, aggregate_type, aggregate_id, payload, customer_id, status, retry_count, locked_at, locked_by, last_error, occurred_at, published_at FROM outbox_events WHERE id = ?"
	var e OutboxEvent
	var corr sql.NullString
	var lockedAt sql.NullTime
	var lockedBy sql.NullString
	var lastErr sql.NullString
	var pubAt sql.NullTime
	var payloadBytes []byte

	err = db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.EventID, &corr, &e.EventType, &e.AggregateType,
		&e.AggregateID, &payloadBytes, &e.CustomerID, &e.Status,
		&e.RetryCount, &lockedAt, &lockedBy, &lastErr, &e.OccurredAt, &pubAt,
	)
	if err != nil {
		return nil, err
	}

	e.Payload = string(payloadBytes)
	if corr.Valid {
		e.CorrelationID = &corr.String
	}
	if lockedAt.Valid {
		e.LockedAt = &lockedAt.Time
	}
	if lockedBy.Valid {
		e.LockedBy = &lockedBy.String
	}
	if lastErr.Valid {
		e.LastError = &lastErr.String
	}
	if pubAt.Valid {
		e.PublishedAt = &pubAt.Time
	}

	return &e, nil
}

// extractIDFromPath parses the ID from path e.g. /api/admin/outbox-events/{id}/retry
func extractIDFromPath(path string) (int64, error) {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "outbox-events" && i+1 < len(parts) {
			idStr := parts[i+1]
			if idx := strings.Index(idStr, "?"); idx != -1 {
				idStr = idStr[:idx]
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err == nil {
				return id, nil
			}
		}
	}
	return 0, fmt.Errorf("id not found in path")
}

// HandleGetOutboxEvents is a reusable Gorilla Mux HTTP handler for querying outbox events.
func HandleGetOutboxEvents(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	customerID := req.Header.Get("X-Customer-Id")
	if customerID == "" {
		customerID = req.URL.Query().Get("customer_id")
	}
	status := req.URL.Query().Get("status")
	eventType := req.URL.Query().Get("event_type")
	eventID := req.URL.Query().Get("event_id")
	correlationID := req.URL.Query().Get("correlation_id")

	limitStr := req.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	pageStr := req.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	events, total, err := QueryOutboxEvents(ctx, db, customerID, status, eventType, eventID, correlationID, page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("database error: %v", err), http.StatusInternalServerError)
		return
	}

	pages := total / limit
	if total%limit > 0 {
		pages++
	}

	resp := PaginatedResponse[OutboxEvent]{
		Data: events,
		Pagination: PaginationMetadata{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: pages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// HandleGetProcessedEvents is a reusable Gorilla Mux HTTP handler for querying processed events.
func HandleGetProcessedEvents(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	eventID := req.URL.Query().Get("event_id")

	limitStr := req.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	pageStr := req.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	events, total, err := QueryProcessedEvents(ctx, db, eventID, page, limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("database error: %v", err), http.StatusInternalServerError)
		return
	}

	pages := total / limit
	if total%limit > 0 {
		pages++
	}

	resp := PaginatedResponse[ProcessedEvent]{
		Data: events,
		Pagination: PaginationMetadata{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: pages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// HandleRetryOutboxEvent is a reusable Gorilla Mux HTTP handler for retrying an outbox event.
func HandleRetryOutboxEvent(db *sql.DB, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id, err := extractIDFromPath(req.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e, err := RetryEvent(ctx, db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(e)
}
