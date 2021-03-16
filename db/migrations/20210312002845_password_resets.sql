-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_resets(
    id BIGSERIAL PRIMARY KEY,
    email text NOT NULL,
    token text NOT NULL UNIQUE,
    valid bool NOT NULL DEFAULT true
);
CREATE INDEX password_resets_idx on password_resets(token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_resets;
-- +goose StatementEnd
