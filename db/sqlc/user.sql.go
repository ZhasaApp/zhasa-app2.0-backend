// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: user.sql

package generated

import (
	"context"
	"time"
)

const addBrandToUser = `-- name: AddBrandToUser :exec
INSERT INTO user_brands (user_id, brand_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`

type AddBrandToUserParams struct {
	UserID  int32 `json:"user_id"`
	BrandID int32 `json:"brand_id"`
}

func (q *Queries) AddBrandToUser(ctx context.Context, arg AddBrandToUserParams) error {
	_, err := q.db.ExecContext(ctx, addBrandToUser, arg.UserID, arg.BrandID)
	return err
}

const addRoleToUser = `-- name: AddRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`

type AddRoleToUserParams struct {
	UserID int32 `json:"user_id"`
	RoleID int32 `json:"role_id"`
}

func (q *Queries) AddRoleToUser(ctx context.Context, arg AddRoleToUserParams) error {
	_, err := q.db.ExecContext(ctx, addRoleToUser, arg.UserID, arg.RoleID)
	return err
}

const addUserToBranch = `-- name: AddUserToBranch :exec
INSERT INTO branch_users (user_id, branch_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`

type AddUserToBranchParams struct {
	UserID   int32 `json:"user_id"`
	BranchID int32 `json:"branch_id"`
}

func (q *Queries) AddUserToBranch(ctx context.Context, arg AddUserToBranchParams) error {
	_, err := q.db.ExecContext(ctx, addUserToBranch, arg.UserID, arg.BranchID)
	return err
}

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
SELECT u.id,
       u.phone,
       u.first_name,
       u.last_name,
       u.avatar_url
FROM user_avatar_view u
WHERE u.phone = $1
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

const getUsersByBranchBrandRole = `-- name: GetUsersByBranchBrandRole :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url
FROM user_avatar_view u
         JOIN user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         JOIN branch_users bu ON u.id = bu.user_id AND bu.branch_id = $2
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $3
`

type GetUsersByBranchBrandRoleParams struct {
	BrandID  int32 `json:"brand_id"`
	BranchID int32 `json:"branch_id"`
	RoleID   int32 `json:"role_id"`
}

type GetUsersByBranchBrandRoleRow struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
}

func (q *Queries) GetUsersByBranchBrandRole(ctx context.Context, arg GetUsersByBranchBrandRoleParams) ([]GetUsersByBranchBrandRoleRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersByBranchBrandRole, arg.BrandID, arg.BranchID, arg.RoleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByBranchBrandRoleRow
	for rows.Next() {
		var i GetUsersByBranchBrandRoleRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
		); err != nil {
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

const getUsersWithoutRoles = `-- name: GetUsersWithoutRoles :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.phone
FROM users u
WHERE
    NOT EXISTS(
        SELECT 1
         FROM user_roles ur
         WHERE ur.user_id = u.id
    ) AND (u.last_name || ' ' || u.first_name) ILIKE $1::text || '%'
ORDER BY u.created_at DESC
LIMIT 10
`

type GetUsersWithoutRolesRow struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

func (q *Queries) GetUsersWithoutRoles(ctx context.Context, search string) ([]GetUsersWithoutRolesRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithoutRoles, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersWithoutRolesRow
	for rows.Next() {
		var i GetUsersWithoutRolesRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Phone,
		); err != nil {
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

const setUserBrandGoal = `-- name: SetUserBrandGoal :exec
INSERT INTO user_brand_sale_type_goals (user_id, brand_id, sale_type_id, value, from_date, to_date)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (user_id, brand_id, sale_type_id, from_date, to_date) DO
UPDATE
    SET value = $4
`

type SetUserBrandGoalParams struct {
	UserID     int32     `json:"user_id"`
	BrandID    int32     `json:"brand_id"`
	SaleTypeID int32     `json:"sale_type_id"`
	Value      int64     `json:"value"`
	FromDate   time.Time `json:"from_date"`
	ToDate     time.Time `json:"to_date"`
}

func (q *Queries) SetUserBrandGoal(ctx context.Context, arg SetUserBrandGoalParams) error {
	_, err := q.db.ExecContext(ctx, setUserBrandGoal,
		arg.UserID,
		arg.BrandID,
		arg.SaleTypeID,
		arg.Value,
		arg.FromDate,
		arg.ToDate,
	)
	return err
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
