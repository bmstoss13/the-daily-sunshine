-- name: SavePost :exec
INSERT INTO saved_posts (user_id, post_id)
VALUES ($1, $2);

-- name: UnsavePost :exec
DELETE FROM saved_posts
WHERE user_id = $1 AND post_id = $2;

-- name: CheckIfPostIsSaved :one
SELECT EXISTS(
    SELECT 1 FROM saved_posts
    WHERE user_id = $1 AND post_id = $2
);

-- name: ListSavedPostsForUser :many
SELECT p.* 
FROM posts p
INNER JOIN saved_posts sp ON p.id = sp.post_id
WHERE sp.user_id = $1 AND p.deleted_at IS NULL
ORDER BY sp.created_at DESC
LIMIT $2 OFFSET $3;