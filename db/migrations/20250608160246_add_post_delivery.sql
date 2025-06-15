-- +goose Up
CREATE TABLE IF NOT EXISTS post_delivery (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL REFERENCES post(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscription(id) ON DELETE CASCADE,
    opened BOOLEAN NOT NULL DEFAULT FALSE
);
-- +goose Down
DROP TABLE IF EXISTS post_delivery;
