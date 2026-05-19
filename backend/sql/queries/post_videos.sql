-- name: SelectPostVideoByID :one
SELECT * FROM post_videos
WHERE id = $1 LIMIT 1;

-- name: SelectPostVideoByPostID :one
SELECT * FROM post_videos
WHERE post_id = $1 LIMIT 1;

-- name: CreatePostVideo :one
INSERT INTO post_videos (
    post_id,
    youtube_video_id,
    video_metadata
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdatePostVideo :one
UPDATE post_videos 
SET
    youtube_video_id = $3,
    video_metadata = $4
WHERE post_videos.id = $1 AND post_id = $2
AND EXISTS (
    SELECT 1 FROM posts
    WHERE posts.id = $2 AND posts.publisher_id = $5
)
RETURNING *;

-- name: DeletePostVideo :one
DELETE FROM post_videos
WHERE post_videos.id = $1 AND post_id = $2
AND EXISTS (
    SELECT 1 FROM posts 
    WHERE posts.id = $2 AND posts.publisher_id = $3
)
RETURNING *;

-- name: ListVideoFeed :many
SELECT
    pv.id AS video_id,
    pv.youtube_video_id,
    pv.video_metadata,
    p.id AS post_id,
    p.title AS post_title,
    p.slug AS post_slug,
    p.publisher_id,
    p.num_rays,
    p.num_comments,
    pv.created_at
FROM post_videos pv
INNER JOIN posts p ON pv.post_id = p.id
WHERE p.deleted_at IS NULL
ORDER BY pv.created_at DESC
LIMIT $1 OFFSET $2;
