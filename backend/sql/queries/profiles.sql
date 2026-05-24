-- name: GetProfileByID :one
SELECT * FROM profiles
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetProfileByUsername :one
SELECT * FROM profiles
WHERE username = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CheckUsernameExists :one
SELECT EXISTS(
    SELECT 1 FROM profiles WHERE username = $1
);

-- name: ListProfiles :many
SELECT * FROM profiles
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProfile :one
INSERT INTO profiles (
    id, 
    first_name, 
    last_name, 
    username, 
    -- subscriber_tier,
    profile_image_url, 
    profile_bio
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET
    first_name = $2,
    last_name = $3,
    username = $4,
    -- subscriber_tier = $5,
    profile_image_url = $5,
    profile_bio = $6,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProfile :one
UPDATE profiles
SET
    deleted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PermanentlyDeleteProfile :one
DELETE FROM profiles
WHERE id = $1
RETURNING *;

-- name: SelectProfilePhoto :one
SELECT profile_image_url FROM profiles
WHERE id = $1 LIMIT 1;

-- name: SetSubscriberTier :one
UPDATE profiles
SET
    subscriber_tier = $2,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetSubscriberTier :one
SELECT subscriber_tier FROM profiles
WHERE id = $1 LIMIT 1;
