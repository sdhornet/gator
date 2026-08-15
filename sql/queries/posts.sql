-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, feed_id, url, title, description, published_at)
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
SELECT posts.*
FROM posts
JOIN feeds ON feeds.id = posts.feed_id
JOIN feed_follows ON feed_follows.feed_id = feeds.id
JOIN users on users.id = feed_follows.user_id
WHERE users.id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
