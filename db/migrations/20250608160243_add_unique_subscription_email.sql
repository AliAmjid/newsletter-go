-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS subscription_unique_newsletter_email ON subscription (newsletter_id, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS subscription_unique_newsletter_email;
-- +goose StatementEnd
