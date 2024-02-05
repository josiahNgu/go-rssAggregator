-- name: CreateFeed :one
INSERT INTO feeds(id,created_at,updated_at,name,url,user_uuid)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetUserFeeds :many
SELECT * FROM feeds WHERE user_uuid = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() where id = $1 RETURNING *;