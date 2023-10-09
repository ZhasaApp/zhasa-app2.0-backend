// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package generated

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (phone, first_name, last_name)
VALUES ($1, $2, $3) ON CONFLICT (phone)
DO
UPDATE SET first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name
    RETURNING id
`

type CreateUserParams struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Phone, arg.FirstName, arg.LastName)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createUserCode = `-- name: CreateUserCode :one
INSERT INTO users_codes(user_id, code)
VALUES ($1, $2) RETURNING id
`

type CreateUserCodeParams struct {
	UserID int32 `json:"user_id"`
	Code   int32 `json:"code"`
}

func (q *Queries) CreateUserCode(ctx context.Context, arg CreateUserCodeParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createUserCode, arg.UserID, arg.Code)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const deleteUserAvatar = `-- name: DeleteUserAvatar :exec
DELETE
FROM users_avatars
WHERE user_id = $1
`

func (q *Queries) DeleteUserAvatar(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteUserAvatar, userID)
	return err
}

const getAuthCodeById = `-- name: GetAuthCodeById :one
SELECT id, user_id, code, created_at
FROM users_codes
WHERE id = $1
`

func (q *Queries) GetAuthCodeById(ctx context.Context, id int32) (UsersCode, error) {
	row := q.db.QueryRowContext(ctx, getAuthCodeById, id)
	var i UsersCode
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Code,
		&i.CreatedAt,
	)
	return i, err
}

const getUserBranch = `-- name: GetUserBranch :one
SELECT b.title, b.id
FROM users u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN branches b ON bu.branch_id = b.id
WHERE u.id = $1
`

type GetUserBranchRow struct {
	Title string `json:"title"`
	ID    int32  `json:"id"`
}

func (q *Queries) GetUserBranch(ctx context.Context, id int32) (GetUserBranchRow, error) {
	row := q.db.QueryRowContext(ctx, getUserBranch, id)
	var i GetUserBranchRow
	err := row.Scan(&i.Title, &i.ID)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT u.id, phone, first_name, last_name, avatar_url, ur.id, user_id, role_id, r.id, title, key, description, created_at
FROM user_avatar_view u
         JOIN user_roles ur on u.id = ur.user_id
         JOIN roles r on ur.role_id = r.id
WHERE u.id = $1
`

type GetUserByIdRow struct {
	ID          int32     `json:"id"`
	Phone       string    `json:"phone"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	AvatarUrl   string    `json:"avatar_url"`
	ID_2        int32     `json:"id_2"`
	UserID      int32     `json:"user_id"`
	RoleID      int32     `json:"role_id"`
	ID_3        int32     `json:"id_3"`
	Title       string    `json:"title"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (q *Queries) GetUserById(ctx context.Context, id int32) (GetUserByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i GetUserByIdRow
	err := row.Scan(
		&i.ID,
		&i.Phone,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.ID_2,
		&i.UserID,
		&i.RoleID,
		&i.ID_3,
		&i.Title,
		&i.Key,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, phone, first_name, last_name, avatar_url
FROM user_avatar_view
WHERE phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone string) (UserAvatarView, error) {
	row := q.db.QueryRowContext(ctx, getUserByPhone, phone)
	var i UserAvatarView
	err := row.Scan(
		&i.ID,
		&i.Phone,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
	)
	return i, err
}

const uploadUserAvatar = `-- name: UploadUserAvatar :exec
INSERT INTO users_avatars(user_id, avatar_url)
VALUES ($1, $2) ON CONFLICT (user_id)
DO
UPDATE SET avatar_url = EXCLUDED.avatar_url
`

type UploadUserAvatarParams struct {
	UserID    int32  `json:"user_id"`
	AvatarUrl string `json:"avatar_url"`
}

func (q *Queries) UploadUserAvatar(ctx context.Context, arg UploadUserAvatarParams) error {
	_, err := q.db.ExecContext(ctx, uploadUserAvatar, arg.UserID, arg.AvatarUrl)
	return err
}
