-- name: SelectPostByID :one
SELECT * FROM posts
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: SelectPostBySlug :one
SELECT * FROM posts
WHERE slug = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListPosts :many
SELECT 
    p.id,
    p.title,
    p.slug,
    p.post_content,
    p.num_rays,
    p.num_comments,
    p.created_at,
    -- Join the Cover Image Data
    pi.image_url AS cover_image_url,
    pi.alt_text AS cover_image_alt,
    -- Join the Publisher Data
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar
FROM posts p
-- Grab the image ONLY if it is marked as the cover image
LEFT JOIN post_images pi ON p.id = pi.post_id AND pi.is_cover_image = true
-- Grab the publisher's profile
LEFT JOIN profiles pr ON p.publisher_id = pr.id
WHERE p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreatePost :one
INSERT INTO posts (
    id,
    publisher_id,
    title,
    slug,
    post_content,
    num_rays,
    num_comments
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
SET 
    title = $3,
    slug = $4,
    post_content = $5,
    num_rays = $6,
    num_comments = $7,
    updated_at = NOW()
WHERE id = $1 AND publisher_id = $2 AND deleted_at IS NULL -- do we need to include publisher id to limit only publishers to update?
RETURNING *;

-- name: SoftDeletePost :one
UPDATE posts
SET
    deleted_at = NOW()
WHERE id = $1 AND publisher_id = $2
RETURNING *;

-- name: PermanentlyDeletePost :one
DELETE FROM posts
WHERE id = $1 AND publisher_id = $2
RETURNING *;

-- name: CountPostsByUserToday :one
SELECT COUNT(*) FROM posts
WHERE publisher_id = $1 
AND created_at >= CURRENT_DATE
AND deleted_at IS NULL;