-- name: GetUsers :many
SELECT * FROM users;


-- name: CreateUser :one
INSERT INTO users (
first_name, last_name,
email, phone_number, role,password  )
VALUES ( $1, $2, $3, $4, $5,$6)
RETURNING *;


-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
);



-- name: GetUserWithEmail :one
SELECT * FROM users WHERE $1=email;
-- name: GetUser :one
SELECT * FROM users WHERE $1=id;


-- name: UpdatePassword :exec
UPDATE users
SET 
  password = $1
WHERE email = $2;