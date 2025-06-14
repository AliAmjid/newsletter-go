-- +goose Up
ALTER TABLE post ALTER COLUMN published_at DROP NOT NULL;
ALTER TABLE post ALTER COLUMN published_at DROP DEFAULT;
-- +goose Down
ALTER TABLE post ALTER COLUMN published_at SET NOT NULL;
ALTER TABLE post ALTER COLUMN published_at SET DEFAULT NOW();
