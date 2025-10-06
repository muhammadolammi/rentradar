
-- name: GetListings :many
SELECT *
FROM listings
WHERE
  (location = coalesce(sqlc.narg('location'), location) OR coalesce(sqlc.narg('location'), location) IS NULL)
AND (price >= coalesce(sqlc.narg('min_price'), min_price) OR coalesce(sqlc.narg('min_price'), min_price) IS NULL)
AND (price<= coalesce(sqlc.narg('max_price'), max_price) OR coalesce(sqlc.narg('max_price'), max_price) IS NULL)
 AND (type = coalesce(sqlc.narg('type'), type) OR coalesce(sqlc.narg('type'), type) IS NULL)
ORDER BY created_at DESC
LIMIT sqlc.arg('limit');