-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invites(
    id BIGSERIAL PRIMARY KEY,
    invited_by text NOT NULL,
    invitee_email text NOT NULL,
    token text NOT NULL UNIQUE,
    valid boolean NOT NULL DEFAULT true,
    created_at timestamp WITH time zone DEFAULT now()
);
CREATE INDEX invites_idx on invites(invited_by, invitee_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invites;
-- +goose StatementEnd
