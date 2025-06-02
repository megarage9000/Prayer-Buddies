-- name: CreatePrayer :one
INSERT INTO prayers(id, created_at, updated_at, sender, receiver, prayer) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPrayersForUser :many
SELECT * FROM prayers
WHERE prayers.sender = $1;