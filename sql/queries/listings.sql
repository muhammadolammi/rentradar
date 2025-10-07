
-- name: GetListings :many
SELECT *
FROM listings
WHERE
  (location = coalesce(sqlc.narg('location'), location))
  AND (price >= coalesce(sqlc.narg('min_price')::bigint, price))
  AND (price <= coalesce(sqlc.narg('max_price')::bigint, price))
  AND (house_type = coalesce(sqlc.narg('house_type'), house_type))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateListing :one
INSERT INTO listings (
agent_id, title,
description, rent_type, price,location,house_type,images, status  )
VALUES ( $1, $2, $3, $4, $5,$6,$7,$8,$9)
RETURNING *;




-- name: GetListing :one
SELECT * FROM listings WHERE $1=id;