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

-- name: GetReceivedPrayersForUser :many
SELECT * FROM prayers
WHERE prayers.receiver = $1
ORDER BY prayers.created_at DESC
LIMIT $2;

-- name: GetSentPrayersFromUser :many
SELECT * FROM prayers
WHERE prayers.sender = $1
ORDER BY prayers.created_at DESC
LIMIT $2;