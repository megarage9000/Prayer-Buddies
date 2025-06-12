-- name: CreateFriendRequest :one
INSERT INTO friends(user_id, friend_id, created_at)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;

-- name: UpdateFriendRequest :exec
UPDATE friends
SET status = $3
WHERE friends.user_id = $1 AND friends.friend_id = $2 AND friend.status != 'pending';

-- name: DenyFriendRequest :exec
DELETE FROM friends
WHERE friends.user_id = $1 AND friends.friend_id = $2;   

-- name: AcceptFriendRequest :one
INSERT INTO friends(user_id, friend_id, status, created_at)
VALUES(
    $1,
    $2,
    'accepted',
    $3
)
RETURNING *;

-- name: GetFriendsFromUser :many
SELECT friends.friend_id FROM friends
WHERE friends.user_id = $1;