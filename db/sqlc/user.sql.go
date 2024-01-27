// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: user.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
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

const addDisabledUser = `-- name: AddDisabledUser :exec
INSERT INTO disabled_users (user_id)
VALUES ($1) ON CONFLICT DO NOTHING
`

func (q *Queries) AddDisabledUser(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, addDisabledUser, userID)
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

const addUserBranch = `-- name: AddUserBranch :exec
INSERT INTO branch_users (user_id, branch_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`

type AddUserBranchParams struct {
	UserID   int32 `json:"user_id"`
	BranchID int32 `json:"branch_id"`
}

func (q *Queries) AddUserBranch(ctx context.Context, arg AddUserBranchParams) error {
	_, err := q.db.ExecContext(ctx, addUserBranch, arg.UserID, arg.BranchID)
	return err
}

const addUserRole = `-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, (SELECT id FROM roles WHERE key = $2::text)) ON CONFLICT DO NOTHING
`

type AddUserRoleParams struct {
	UserID  int32  `json:"user_id"`
	RoleKey string `json:"role_key"`
}

func (q *Queries) AddUserRole(ctx context.Context, arg AddUserRoleParams) error {
	_, err := q.db.ExecContext(ctx, addUserRole, arg.UserID, arg.RoleKey)
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

const getFilteredUsersWithBranchRolesBrands = `-- name: GetFilteredUsersWithBranchRolesBrands :many
WITH Counted AS (
    SELECT u.id,
           u.first_name,
           u.last_name,
           u.phone,
           r.key                      AS role,
           b.title                    AS branch_title,
           STRING_AGG(bs.title, ', ') AS brands,
           COUNT(*) OVER()            AS total_count,
           CASE
               WHEN du.user_id IS NULL THEN true
               ELSE false
               END                        AS is_active
    FROM users u
             JOIN user_roles ur ON u.id = ur.user_id
             JOIN roles r ON ur.role_id = r.id
             LEFT JOIN branch_users bu ON u.id = bu.user_id
             LEFT JOIN user_brands ub ON u.id = ub.user_id
             LEFT JOIN brands bs ON ub.brand_id = bs.id
             LEFT JOIN branches b ON bu.branch_id = b.id
             LEFT JOIN disabled_users du ON u.id = du.user_id
    WHERE (last_name || ' ' || first_name) ILIKE '%' || $5::text || '%'
      AND ($6::text[] IS NULL OR r.key = ANY($6))
      AND ($7::int[] IS NULL OR bs.id = ANY($7))
      AND ($8::int[] IS NULL OR b.id = ANY($8))
      AND (du.user_id IS NULL)
    GROUP BY u.id, u.first_name, u.last_name, b.title, du.user_id, r.key
)
SELECT id,
       first_name,
       last_name,
       phone,
       role,
       branch_title,
       brands,
       total_count,
       is_active
FROM Counted
ORDER BY
    CASE WHEN $3::text = 'fio' AND $4::text = 'asc' THEN first_name END ASC,
    CASE WHEN $3 = 'fio' AND $4 = 'asc' THEN last_name END ASC,
    CASE WHEN $3 = 'fio' AND $4 = 'desc' THEN first_name END DESC,
    CASE WHEN $3 = 'fio' AND $4 = 'desc' THEN last_name END DESC,
    CASE WHEN $3 = 'branch' AND $4 = 'asc' THEN branch_title END ASC,
    CASE WHEN $3 = 'branch' AND $4 = 'desc' THEN branch_title END DESC,
    first_name, last_name, id DESC
LIMIT $1 OFFSET $2
`

type GetFilteredUsersWithBranchRolesBrandsParams struct {
	Limit     int32    `json:"limit"`
	Offset    int32    `json:"offset"`
	SortField string   `json:"sort_field"`
	SortType  string   `json:"sort_type"`
	Search    string   `json:"search"`
	RoleKeys  []string `json:"role_keys"`
	BrandIds  []int32  `json:"brand_ids"`
	BranchIds []int32  `json:"branch_ids"`
}

type GetFilteredUsersWithBranchRolesBrandsRow struct {
	ID          int32          `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Phone       string         `json:"phone"`
	Role        string         `json:"role"`
	BranchTitle sql.NullString `json:"branch_title"`
	Brands      []byte         `json:"brands"`
	TotalCount  int64          `json:"total_count"`
	IsActive    bool           `json:"is_active"`
}

func (q *Queries) GetFilteredUsersWithBranchRolesBrands(ctx context.Context, arg GetFilteredUsersWithBranchRolesBrandsParams) ([]GetFilteredUsersWithBranchRolesBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFilteredUsersWithBranchRolesBrands,
		arg.Limit,
		arg.Offset,
		arg.SortField,
		arg.SortType,
		arg.Search,
		pq.Array(arg.RoleKeys),
		pq.Array(arg.BrandIds),
		pq.Array(arg.BranchIds),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFilteredUsersWithBranchRolesBrandsRow
	for rows.Next() {
		var i GetFilteredUsersWithBranchRolesBrandsRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Phone,
			&i.Role,
			&i.BranchTitle,
			&i.Brands,
			&i.TotalCount,
			&i.IsActive,
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

const getUsersWithBranchBrands = `-- name: GetUsersWithBranchBrands :many
WITH Counted AS (
    SELECT u.id,
           u.first_name,
           u.last_name,
           u.phone,
           b.title                    AS branch_title,
           STRING_AGG(bs.title, ', ') AS brands,
           COUNT(*) OVER()            AS total_count
    FROM users u
             JOIN branch_users bu ON u.id = bu.user_id
             JOIN user_brands ub ON u.id = ub.user_id
             JOIN brands bs ON ub.brand_id = bs.id
             JOIN branches b ON bu.branch_id = b.id
    WHERE (last_name || ' ' || first_name) ILIKE '%' || $3::text || '%'
    GROUP BY u.id, u.first_name, u.last_name, b.title
)
SELECT id,
       first_name,
       last_name,
       phone,
       branch_title,
       brands,
       total_count
FROM Counted
ORDER BY first_name, last_name, id DESC
LIMIT $1 OFFSET $2
`

type GetUsersWithBranchBrandsParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Search string `json:"search"`
}

type GetUsersWithBranchBrandsRow struct {
	ID          int32  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	BranchTitle string `json:"branch_title"`
	Brands      []byte `json:"brands"`
	TotalCount  int64  `json:"total_count"`
}

func (q *Queries) GetUsersWithBranchBrands(ctx context.Context, arg GetUsersWithBranchBrandsParams) ([]GetUsersWithBranchBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithBranchBrands, arg.Limit, arg.Offset, arg.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersWithBranchBrandsRow
	for rows.Next() {
		var i GetUsersWithBranchBrandsRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Phone,
			&i.BranchTitle,
			&i.Brands,
			&i.TotalCount,
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

const getUsersWithBranchRolesBrands = `-- name: GetUsersWithBranchRolesBrands :many
WITH Counted AS (
    SELECT u.id,
           u.first_name,
           u.last_name,
           u.phone,
           b.title                    AS branch_title,
           STRING_AGG(bs.title, ', ') AS brands,
           COUNT(*) OVER()            AS total_count,
           CASE
               WHEN du.user_id IS NULL THEN true
               ELSE false
           END                        AS is_active
    FROM users u
             JOIN user_roles ur ON u.id = ur.user_id
             JOIN roles r ON ur.role_id = r.id AND r.key = $1
             JOIN branch_users bu ON u.id = bu.user_id
             JOIN user_brands ub ON u.id = ub.user_id
             JOIN brands bs ON ub.brand_id = bs.id
             JOIN branches b ON bu.branch_id = b.id
             LEFT JOIN disabled_users du ON u.id = du.user_id
    WHERE (last_name || ' ' || first_name) ILIKE '%' || $4::text || '%'
    GROUP BY u.id, u.first_name, u.last_name, b.title, du.user_id
)
SELECT id,
       first_name,
       last_name,
       phone,
       branch_title,
       brands,
       total_count,
       is_active
FROM Counted
ORDER BY first_name, last_name, id DESC
LIMIT $2 OFFSET $3
`

type GetUsersWithBranchRolesBrandsParams struct {
	Key    string `json:"key"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Search string `json:"search"`
}

type GetUsersWithBranchRolesBrandsRow struct {
	ID          int32  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	BranchTitle string `json:"branch_title"`
	Brands      []byte `json:"brands"`
	TotalCount  int64  `json:"total_count"`
	IsActive    bool   `json:"is_active"`
}

func (q *Queries) GetUsersWithBranchRolesBrands(ctx context.Context, arg GetUsersWithBranchRolesBrandsParams) ([]GetUsersWithBranchRolesBrandsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithBranchRolesBrands,
		arg.Key,
		arg.Limit,
		arg.Offset,
		arg.Search,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersWithBranchRolesBrandsRow
	for rows.Next() {
		var i GetUsersWithBranchRolesBrandsRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Phone,
			&i.BranchTitle,
			&i.Brands,
			&i.TotalCount,
			&i.IsActive,
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
       u.last_name
FROM users u
    LEFT JOIN user_roles ur ON u.id = ur.user_id
WHERE ur.user_id IS NULL AND (u.last_name || ' ' || u.first_name) ILIKE '%' || $1::text || '%'
ORDER BY u.created_at DESC
LIMIT 25
`

type GetUsersWithoutRolesRow struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
		if err := rows.Scan(&i.ID, &i.FirstName, &i.LastName); err != nil {
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

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET first_name = $1, last_name = $2, phone = $3
WHERE id = $4
`

type UpdateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	ID        int32  `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
		arg.ID,
	)
	return err
}

const updateUserBranch = `-- name: UpdateUserBranch :exec
UPDATE branch_users
SET branch_id = $1
WHERE user_id = $2
`

type UpdateUserBranchParams struct {
	BranchID int32 `json:"branch_id"`
	UserID   int32 `json:"user_id"`
}

func (q *Queries) UpdateUserBranch(ctx context.Context, arg UpdateUserBranchParams) error {
	_, err := q.db.ExecContext(ctx, updateUserBranch, arg.BranchID, arg.UserID)
	return err
}

const updateUserRole = `-- name: UpdateUserRole :exec
UPDATE user_roles
SET role_id = (SELECT id FROM roles WHERE key = $2::text)
WHERE user_id = $1
`

type UpdateUserRoleParams struct {
	UserID  int32  `json:"user_id"`
	RoleKey string `json:"role_key"`
}

func (q *Queries) UpdateUserRole(ctx context.Context, arg UpdateUserRoleParams) error {
	_, err := q.db.ExecContext(ctx, updateUserRole, arg.UserID, arg.RoleKey)
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
