-- name: DeleteBranchUserByUserId :exec
DELETE FROM branch_users WHERE user_id = $1;
