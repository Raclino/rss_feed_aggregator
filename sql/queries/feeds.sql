-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES 
    ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT
    f.name,
    f.url,
    u.name
FROM
    feeds f
    INNER JOIN (
        SELECT
            id,
            name
        FROM
            users
    ) u on f.user_id = u.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES
        ($1, $2, $3, $4, $5)
        RETURNING *
)
SELECT
    isf.id,
    isf.created_at,
    isf.updated_at,
    isf.user_id,
    isf.feed_id,
    u.name AS user_name,
    f.name AS feed_name
FROM
    inserted_feed_follow isf
    INNER JOIN users u ON isf.user_id = u.id
    INNER JOIN feeds f ON isf.feed_id = f.id;

-- name: GetFeedByUrl :one
SELECT 
    f.name,
    f.id,
    f.url,
    f.user_id
FROM 
    feeds f
WHERE f.url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    *
FROM feed_follows ff
    INNER JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1;
