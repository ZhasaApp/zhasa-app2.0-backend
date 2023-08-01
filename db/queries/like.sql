-- name: AddLike :one
INSERT INTO likes(user_id, post_id)
    VALUES ($1, $2)
RETURNING *;

-- name: GetUserPostLike :one
SELECT user_id FROM likes
WHERE user_id = $1 AND post_id = $2;

-- name: GetPostLikesCount :one
SELECT COUNT(user_id) FROM likes
WHERE post_id = $1;

-- name: GetPostLikedUsers :many
SELECT l.user_id, u.first_name, u.last_name FROM likes l
JOIN users u
ON l.user_id = u.id
WHERE l.post_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteLike :exec
DELETE FROM likes
WHERE user_id = $1 AND post_id = $2;
