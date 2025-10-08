

-- name: CreateAlert :one
INSERT INTO alerts (
user_id, min_price,max_price, location, property_type_id,contact_method )
VALUES ( $1, $2, $3, $4, $5,$6)
RETURNING *;


-- name: GetAlert :one
SELECT * FROM alerts WHERE $1=id;


-- name: GetUserAlerts :many
SELECT * FROM alerts WHERE $1=user_id;