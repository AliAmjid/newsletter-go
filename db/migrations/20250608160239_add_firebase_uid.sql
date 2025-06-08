-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS firebase_uid TEXT UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" DROP COLUMN IF EXISTS firebase_uid;
-- +goose StatementEnd
