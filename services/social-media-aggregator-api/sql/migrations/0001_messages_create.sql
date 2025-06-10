-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    modified_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    external_id text NOT NULL UNIQUE,
    author text NOT NULL,
    body text NOT NULL,
    platform text NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS messages;
