
-- name: RegisterUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE users.email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE users.username = $1;