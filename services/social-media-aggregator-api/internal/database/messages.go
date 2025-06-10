package database

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID         int64     `json:"id"`
	ExternalId string    `json:"external_id"`
	Author     string    `json:"author"`
	Body       string    `json:"body"`
	Platform   string    `json:"platform"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
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
	body,
	platform
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, modified_at`

	args := []any{
		message.ExternalId,
		message.Author,
		message.Body,
		message.Platform,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&message.ID,
		&message.CreatedAt,
		&message.ModifiedAt,
	)

}

func (m *MessageModel) GetList(limit, offset int) ([]*Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT id, external_id, author, body, platform, created_at, modified_at  FROM messages
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	args := []any{limit, offset}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID,
			&message.ExternalId,
			&message.Author,
			&message.Body,
			&message.Platform,
			&message.CreatedAt,
			&message.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil

}
