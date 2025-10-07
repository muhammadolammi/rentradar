-- name: GetAgents :many
SELECT * FROM agents;


-- name: CreateAgent :one
INSERT INTO agents (
user_id, company_name)
VALUES ( $1, $2)
RETURNING *; 


-- name: GetAgentWithUserId :one
SELECT * FROM agents WHERE $1=user_id;
-- name: GetAgent :one
SELECT * FROM agents WHERE $1=id;

