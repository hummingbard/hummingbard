-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS two_factor_auth(
    id BIGSERIAL PRIMARY KEY,
    user_id text NOT NULL UNIQUE,
    enabled boolean NOT NULL DEFAULT false,
    secret text NOT NULL UNIQUE,
    recovery_codes jsonb
);
CREATE INDEX two_factor_auth_idx on two_factor_auth(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS two_factor_auth;
-- +goose StatementEnd
