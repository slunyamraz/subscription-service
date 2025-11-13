package repository

import (
	"database/sql"
	"encoding/json"
	"s/domain"
)

type OutboxRepo struct {
	db *sql.DB
}

func NewOutbox(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{db}
}

func (r *OutboxRepo) CreateEvent(tx *sql.Tx, event domain.OutboxEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO outbox_events (aggregate_type, aggregate_id, event_type, payload) VALUES ($1, $2, $3, $4)",
		event.AggregateType, event.AggregateID, event.EventType, payload,
	)
	return err
}
