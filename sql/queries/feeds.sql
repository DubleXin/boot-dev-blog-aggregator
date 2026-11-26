-- name: CreateFeed :one
INSERT INTO
    feeds (created_at, updated_at, name, url, user_id)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetFeedUsersList :many
SELECT
    feeds.name,
    feeds.url,
    users.name AS username
FROM
    feeds
    INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT
    *
FROM
    feeds
WHERE
    url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET
    updated_at = $2,
    last_fetched_at = $2
WHERE
    id = $1;

-- name: GetNextFeedToFetch :one
SELECT
    *
FROM
    feeds
WHERE
    user_id = $1
ORDER BY
    last_fetched_at NULLS FIRST
LIMIT
    1;