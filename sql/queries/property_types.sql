
-- name: GetPropertyTypeWithName :one
SELECT * FROM property_types WHERE $1=name;

-- name: GetPropertyTypes :many
SELECT * FROM property_types;

-- name: CreatePropertyType :one
INSERT INTO property_types (
name  )
VALUES ( $1)
RETURNING *;

