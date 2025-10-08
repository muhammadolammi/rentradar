
-- name: GetListings :many
SELECT *
FROM listings
WHERE
  (location = coalesce(sqlc.narg('location'), location))
  AND (price >= coalesce(sqlc.narg('min_price')::bigint, price))
  AND (price <= coalesce(sqlc.narg('max_price')::bigint, price))
  AND (property_type_id = coalesce(sqlc.narg('property_type_id'), property_type_id))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: CreateListing :one
INSERT INTO listings (
agent_id, title,
description, price,location,property_type_id,images, status  )
VALUES ( $1, $2, $3, $4, $5,$6,$7,$8)
RETURNING *;




-- name: GetListing :one
SELECT * FROM listings WHERE $1=id;