-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- this make that a user can't have 2 same row of data for the same feed follow
    UNIQUE(user_id,feed_id)
);
-- +goose Down
DROP TABLE feed_follows;