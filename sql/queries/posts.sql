-- name: CreatePost :one
INSERT INTO
    posts (
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO NOTHING
RETURNING
    *;

-- name: GetPostsByFeed :many
SELECT
    *
FROM
    posts
WHERE
    feed_id = $1
ORDER BY
    published_at DESC NULLS LAST,
    id DESC
LIMIT
    $2;

-- name: GetPostsForUser :many
SELECT
    posts.*
FROM
    posts
    JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE
    feed_follows.user_id = $1
ORDER BY
    posts.published_at DESC NULLS LAST,
    posts.id DESC
LIMIT
    $2;