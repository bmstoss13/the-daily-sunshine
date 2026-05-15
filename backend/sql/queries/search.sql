-- name: SearchProfiles :many
SELECT
    id,
    username,
    first_name,
    last_name,
    profile_image_url,
    subscriber_tier
FROM profiles 
WHERE deleted_at IS NULL
AND (
    username ILIKE '%' || $1::text || '%'
    OR first_name ILIKE '%' || $1::text || '%'
    OR last_name ILIKE '%' || $1::text || '%'
)

ORDER BY
    (username ILIKE $1::text || '%') DESC,
    (first_name ILIKE $1::text || '%') DESC,
    num_rays_received DESC
LIMIT $2 OFFSET $3;

-- name: SearchPosts :many
SELECT
    p.id,
    p.title,
    p.subtitle,
    p.slug,
    p.created_at,
    pr.username AS publisher_username,
    pr.profile_image_url AS publisher_avatar
FROM posts p
LEFT JOIN profiles pr ON p.publisher_id = pr.id
WHERE p.deleted_at IS NULL AND p.status = 'published'
--Combine Title and Subtitle into one and search
AND to_tsvector('english', p.title || ' ' || COALESCE(p.subtitle, '')) @@ websearch_to_tsquery('english', $1)
-- sort by relevance rank (how well matched)
ORDER BY ts_rank(to_tsvector('english', p.title || ' ' || COALESCE(p.subtitle, '')), websearch_to_tsquery('english', $1)) DESC
LIMIT $2 OFFSET $3;
