-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_feed_follow UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;