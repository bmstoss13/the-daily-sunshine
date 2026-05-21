-- name: SelectPostByID :one
SELECT 
    p.*,
    -- Publisher Data (Added First/Last Name)
    pr.first_name AS publisher_first_name,
    pr.last_name AS publisher_last_name,
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar,
    pr.subscriber_tier AS publisher_role,
    -- Image Data
    pi.id AS image_id,
    pi.image_url AS cover_image_url,
    pi.alt_text AS cover_image_alt,
    -- Video Data 
    pv.id AS video_id,
    pv.youtube_video_id,
    pv.video_metadata
FROM posts p
LEFT JOIN profiles pr ON p.publisher_id = pr.id 
LEFT JOIN post_images pi ON p.id = pi.post_id AND pi.is_cover_image = true
LEFT JOIN post_videos pv ON p.id = pv.post_id
WHERE p.id = $1 AND p.deleted_at IS NULL LIMIT 1;

-- name: SelectPostBySlug :one
SELECT 
    p.*,
    -- Publisher Data
    pr.first_name AS publisher_first_name,
    pr.last_name AS publisher_last_name,
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar,
    pr.subscriber_tier AS publisher_role,
    -- Image Data
    pi.id AS image_id,
    pi.image_url AS cover_image_url,
    pi.alt_text AS cover_image_alt,
    -- Video Data
    pv.id AS video_id,
    pv.youtube_video_id,
    pv.video_metadata
FROM posts p
LEFT JOIN profiles pr ON p.publisher_id = pr.id 
LEFT JOIN post_images pi ON p.id = pi.post_id AND pi.is_cover_image = true
LEFT JOIN post_videos pv ON p.id = pv.post_id
WHERE p.slug = $1 AND p.deleted_at IS NULL LIMIT 1;

-- name: ListPosts :many
SELECT 
    p.id,
    p.publisher_id, 
    p.title,
    p.subtitle,
    p.slug,
    p.post_content,
    p.status,       
    p.num_rays,
    p.num_comments,
    p.created_at,
    pi.image_url AS cover_image_url,
    pi.alt_text AS cover_image_alt,
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar
FROM posts p
LEFT JOIN post_images pi ON p.id = pi.post_id AND pi.is_cover_image = true
LEFT JOIN profiles pr ON p.publisher_id = pr.id
-- ONLY SHOW PUBLISHED POSTS in the main feed
WHERE p.deleted_at IS NULL AND p.status = 'published'
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreatePost :one
INSERT INTO posts (
    publisher_id, title, subtitle, slug, post_content, status, num_rays, num_comments
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdatePost :one
UPDATE posts
SET 
    title = $3,
    subtitle = $4,
    slug = $5,
    post_content = $6,
    status = $7,
    num_rays = $8,
    num_comments = $9,
    updated_at = NOW()
WHERE id = $1 AND publisher_id = $2 AND deleted_at IS NULL 
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

-- name: CheckSlugExists :one
SELECT EXISTS(
    SELECT 1 FROM posts WHERE slug = $1
);

--To update later to include videos
-- name: SelectPostsOfTheDay :many
SELECT 
    p.id,
    p.publisher_id, 
    p.title,
    p.subtitle,
    p.slug,
    p.post_content,
    p.status,       
    p.num_rays,
    p.num_comments,
    p.created_at,
    pi.image_url AS cover_image_url,
    pi.alt_text AS cover_image_alt,
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar
FROM posts p
LEFT JOIN post_images pi ON p.id = pi.post_id AND pi.is_cover_image = true
LEFT JOIN profiles pr ON p.publisher_id = pr.id
WHERE p.deleted_at IS NULL AND p.status = 'published'
AND p.created_at >= $1
AND p.created_at < $2
ORDER BY p.num_rays DESC, p.created_at DESC
LIMIT 10;
