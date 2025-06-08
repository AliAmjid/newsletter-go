-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" DROP COLUMN IF EXISTS password_hash;
-- +goose StatementEnd
