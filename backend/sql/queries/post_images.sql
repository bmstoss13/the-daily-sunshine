-- name: SelectPostImageByID :one
SELECT * FROM post_images
WHERE id = $1 LIMIT 1;

-- name: SelectPostImagesByPost :many
SELECT * FROM post_images
WHERE post_id = $1;

-- name: CreatePostImage :one
INSERT INTO post_images (
    id,
    post_id,
    image_url,
    image_description,
    alt_text,
    is_cover_image
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdatePostImage :one
UPDATE post_images
SET
    image_url = $3,
    image_description = $4,
    alt_text = $5,
    is_cover_image = $6
WHERE id = $1 AND post_id = $2
RETURNING *;

-- name: DeletePostImage :one
DELETE FROM post_images
WHERE id = $1 AND post_id = $2
RETURNING *;
