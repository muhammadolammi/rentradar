

-- name: GetUnsentNotifications :many
SELECT * FROM notifications
WHERE status="pending";


-- name: CreateNotification :one
INSERT INTO notifications (
user_id,
sent_at, contact, contact_method, status  )
VALUES ( $1, $2, $3, $4, $5)
RETURNING *;
