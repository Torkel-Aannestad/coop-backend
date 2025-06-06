package database

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID                int64     `json:"id"`
	ExternalId        string    `json:"external_id"`
	Author            string    `json:"author"`
	Title             string    `json:"title"`
	Body              string    `json:"body"`
	Version           int32     `json:"version"`
	ExternalCreatedAt time.Time `json:"external_created_at"`
	CreatedAt         time.Time `json:"-"`
	ModifiedAt        time.Time `json:"-"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(message *Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
	INSERT INTO messages (
	external_id,
	author,
	title,
	body,
	external_created_at
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, modified_at, version`

	args := []any{
		message.ExternalId,
		message.Author,
		message.Title,
		message.Body,
		message.ExternalCreatedAt,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&message.ID,
		&message.CreatedAt,
		&message.ModifiedAt,
		&message.Version,
	)

}
