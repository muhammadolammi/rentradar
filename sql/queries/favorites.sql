-- name: CreateFavorite :one
INSERT INTO favorites (user_id, listing_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserFavorites :many
SELECT * FROM favorites WHERE user_id = $1;
