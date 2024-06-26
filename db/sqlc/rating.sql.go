// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rating.sql

package generated

import (
	"context"
	"time"
)

const getBranchUsersOrderedByRatioForGivenBrand = `-- name: GetBranchUsersOrderedByRatioForGivenBrand :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(r.ratio, 0) AS ratio,
       b.title              AS branch_title,
       b.id                 AS branch_id
FROM user_avatar_view u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN
     branches b ON bu.branch_id = b.id
         JOIN
     user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         JOIN
     user_brand_ratio r ON u.id = r.user_id AND r.from_date = $2 AND r.to_date = $3
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $7
         LEFT JOIN disabled_users du ON u.id = du.user_id
WHERE (r.brand_id = $1 OR r.brand_id IS NULL)
  AND du.user_id IS NULL
  AND b.id = $6
ORDER BY r.ratio DESC
OFFSET $4 LIMIT $5
`

type GetBranchUsersOrderedByRatioForGivenBrandParams struct {
	BrandID  int32     `json:"brand_id"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Offset   int32     `json:"offset"`
	Limit    int32     `json:"limit"`
	ID       int32     `json:"id"`
	RoleID   int32     `json:"role_id"`
}

type GetBranchUsersOrderedByRatioForGivenBrandRow struct {
	ID          int32   `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	AvatarUrl   string  `json:"avatar_url"`
	Ratio       float32 `json:"ratio"`
	BranchTitle string  `json:"branch_title"`
	BranchID    int32   `json:"branch_id"`
}

// SELECT distinct users for given brand ordered by ratio and limited by offset and limit and if there is no any user with ratio let ratio be 0
func (q *Queries) GetBranchUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetBranchUsersOrderedByRatioForGivenBrandParams) ([]GetBranchUsersOrderedByRatioForGivenBrandRow, error) {
	rows, err := q.db.QueryContext(ctx, getBranchUsersOrderedByRatioForGivenBrand,
		arg.BrandID,
		arg.FromDate,
		arg.ToDate,
		arg.Offset,
		arg.Limit,
		arg.ID,
		arg.RoleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBranchUsersOrderedByRatioForGivenBrandRow
	for rows.Next() {
		var i GetBranchUsersOrderedByRatioForGivenBrandRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.Ratio,
			&i.BranchTitle,
			&i.BranchID,
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

const getUsersOrderedByRatioForGivenBrand = `-- name: GetUsersOrderedByRatioForGivenBrand :many
SELECT u.id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       COALESCE(r.ratio, 0) AS ratio,
       b.title              AS branch_title,
       b.id                 AS branch_id
FROM user_avatar_view u
         JOIN
     branch_users bu ON u.id = bu.user_id
         JOIN
     branches b ON bu.branch_id = b.id
         JOIN
     user_brands ub ON u.id = ub.user_id AND ub.brand_id = $1
         JOIN
     user_brand_ratio r ON u.id = r.user_id AND r.from_date = $2 AND r.to_date = $3
         JOIN user_roles ur ON u.id = ur.user_id AND ur.role_id = $6
         LEFT JOIN disabled_users du ON u.id = du.user_id
WHERE (r.brand_id = $1 OR r.brand_id IS NULL)
  AND du.user_id IS NULL
ORDER BY r.ratio DESC
OFFSET $4 LIMIT $5
`

type GetUsersOrderedByRatioForGivenBrandParams struct {
	BrandID  int32     `json:"brand_id"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Offset   int32     `json:"offset"`
	Limit    int32     `json:"limit"`
	RoleID   int32     `json:"role_id"`
}

type GetUsersOrderedByRatioForGivenBrandRow struct {
	ID          int32   `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	AvatarUrl   string  `json:"avatar_url"`
	Ratio       float32 `json:"ratio"`
	BranchTitle string  `json:"branch_title"`
	BranchID    int32   `json:"branch_id"`
}

// SELECT distinct users for given brand ordered by ratio and limited by offset and limit and if there is no any user with ratio let ratio be 0
func (q *Queries) GetUsersOrderedByRatioForGivenBrand(ctx context.Context, arg GetUsersOrderedByRatioForGivenBrandParams) ([]GetUsersOrderedByRatioForGivenBrandRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersOrderedByRatioForGivenBrand,
		arg.BrandID,
		arg.FromDate,
		arg.ToDate,
		arg.Offset,
		arg.Limit,
		arg.RoleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersOrderedByRatioForGivenBrandRow
	for rows.Next() {
		var i GetUsersOrderedByRatioForGivenBrandRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.Ratio,
			&i.BranchTitle,
			&i.BranchID,
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
