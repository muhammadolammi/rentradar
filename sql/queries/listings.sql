
-- name: GetListings :many
-- SELECT *
-- FROM listings
-- WHERE
--   (location = coalesce(sqlc.narg('location'), location) OR coalesce(sqlc.narg('location'), location) IS NULL)
-- AND (price >= coalesce(sqlc.narg('min_price'), min_price::numeric) OR coalesce(sqlc.narg('min_price'), min_price::numeric) IS NULL)
-- AND (price<= coalesce(sqlc.narg('max_price')::numeric, max_price) OR coalesce(sqlc.narg('max_price')::numeric, max_price) IS NULL)
--  AND (type = coalesce(sqlc.narg('type'), type) OR coalesce(sqlc.narg('type'), type) IS NULL)
-- ORDER BY created_at DESC
-- LIMIT sqlc.arg('limit')
-- OFFSET sqlc.arg('offset');

SELECT *
FROM listings
WHERE
  (location = coalesce(sqlc.narg('location'), location))
  AND (price >= coalesce(sqlc.narg('min_price')::bigint, price))
  AND (price <= coalesce(sqlc.narg('max_price')::bigint, price))
  AND (house_type = coalesce(sqlc.narg('type'), type))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateListing :one
-- INSERT INTO listings (
-- agent_id, title,
-- description, rent_type, price,location, latitude,longtitude,type,images, status  )
-- VALUES ( $1, $2, $3, $4, $5,$6,$7,$8,$9,$10,$11)
-- RETURNING *;
INSERT INTO listings (
agent_id, title,
description, rent_type, price,location,house_type,images, status  )
VALUES ( $1, $2, $3, $4, $5,$6,$7,$8,$9)
RETURNING *;