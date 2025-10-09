

-- -- name: CreateAlert :one
-- INSERT INTO alerts (
-- user_id, min_price,max_price, location, property_type,contact_method )
-- VALUES ( $1, $2, $3, $4, $5,$6)
-- RETURNING *;


-- name: GetAlert :one
SELECT * FROM alerts WHERE $1=id;





-- name: GetAlertsForListing :many
SELECT *
FROM alerts
WHERE
  (location = $1  OR location = 'any')
  AND (property_type = $2  OR property_type = 'any')
  AND (min_price IS NULL OR min_price <= $3)
  AND (max_price IS NULL OR max_price >= $3);




-- -- name: GetUserAlerts :many
-- SELECT * FROM alerts WHERE $1=user_id;