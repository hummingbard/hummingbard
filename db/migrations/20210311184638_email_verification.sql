-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS email_verification(
    id BIGSERIAL PRIMARY KEY,
    email text NOT NULL,
    token text NOT NULL UNIQUE,
    valid bool NOT NULL DEFAULT true
);
CREATE INDEX email_verification_idx on email_verification(token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_verification;
-- +goose StatementEnd
