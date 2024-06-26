// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: like.sql

package generated

import (
	"context"
)

const addLike = `-- name: AddLike :one
INSERT INTO likes(user_id, post_id)
    VALUES ($1, $2)
RETURNING user_id, post_id
`

type AddLikeParams struct {
	UserID int32 `json:"user_id"`
	PostID int32 `json:"post_id"`
}

func (q *Queries) AddLike(ctx context.Context, arg AddLikeParams) (Like, error) {
	row := q.db.QueryRowContext(ctx, addLike, arg.UserID, arg.PostID)
	var i Like
	err := row.Scan(&i.UserID, &i.PostID)
	return i, err
}

const deleteLike = `-- name: DeleteLike :exec
DELETE FROM likes
WHERE user_id = $1 AND post_id = $2
`

type DeleteLikeParams struct {
	UserID int32 `json:"user_id"`
	PostID int32 `json:"post_id"`
}

func (q *Queries) DeleteLike(ctx context.Context, arg DeleteLikeParams) error {
	_, err := q.db.ExecContext(ctx, deleteLike, arg.UserID, arg.PostID)
	return err
}

const getPostLikedUsers = `-- name: GetPostLikedUsers :many
SELECT l.user_id, u.first_name, u.last_name FROM likes l
JOIN users u
ON l.user_id = u.id
WHERE l.post_id = $1
LIMIT $2
OFFSET $3
`

type GetPostLikedUsersParams struct {
	PostID int32 `json:"post_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPostLikedUsersRow struct {
	UserID    int32  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (q *Queries) GetPostLikedUsers(ctx context.Context, arg GetPostLikedUsersParams) ([]GetPostLikedUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostLikedUsers, arg.PostID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostLikedUsersRow
	for rows.Next() {
		var i GetPostLikedUsersRow
		if err := rows.Scan(&i.UserID, &i.FirstName, &i.LastName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostLikesCount = `-- name: GetPostLikesCount :one
SELECT COUNT(user_id) FROM likes
WHERE post_id = $1
`

func (q *Queries) GetPostLikesCount(ctx context.Context, postID int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPostLikesCount, postID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserPostLike = `-- name: GetUserPostLike :one
SELECT user_id FROM likes
WHERE user_id = $1 AND post_id = $2
`

type GetUserPostLikeParams struct {
	UserID int32 `json:"user_id"`
	PostID int32 `json:"post_id"`
}

func (q *Queries) GetUserPostLike(ctx context.Context, arg GetUserPostLikeParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, getUserPostLike, arg.UserID, arg.PostID)
	var user_id int32
	err := row.Scan(&user_id)
	return user_id, err
}
