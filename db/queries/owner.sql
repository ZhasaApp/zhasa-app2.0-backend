

-- name: GetOwnerByUserId :one
SELECT *
FROM owners_view o
WHERE o.user_id = $1;