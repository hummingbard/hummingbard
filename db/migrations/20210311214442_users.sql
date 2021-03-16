-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    user_id text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    email_verified boolean NOT NULL default false,
    access_token text NOT NULL,
    created_at timestamp WITH time zone DEFAULT now()
);
CREATE INDEX users_idx on users(user_id, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
