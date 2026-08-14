-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, feed_id, url, title, description, publised_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, 