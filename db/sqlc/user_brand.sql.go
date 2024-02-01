// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user_brand.sql

package generated

import (
	"context"
	"time"
)

const deleteUserBrandByUserId = `-- name: DeleteUserBrandByUserId :exec
DELETE FROM user_brands
WHERE user_id = $1
`

func (q *Queries) DeleteUserBrandByUserId(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteUserBrandByUserId, userID)
	return err
}

const getUserBrand = `-- name: GetUserBrand :one
SELECT ub.id AS user_brand
FROM user_brands ub
WHERE ub.user_id = $1
  AND ub.brand_id = $2
`

type GetUserBrandParams struct {
	UserID  int32 `json:"user_id"`
	BrandID int32 `json:"brand_id"`
}

func (q *Queries) GetUserBrand(ctx context.Context, arg GetUserBrandParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, getUserBrand, arg.UserID, arg.BrandID)
	var user_brand int32
	err := row.Scan(&user_brand)
	return user_brand, err
}

const getUserBrandGoal = `-- name: GetUserBrandGoal :one
SELECT COALESCE(goals.value, 0)
FROM user_brand_sale_type_goals goals
WHERE goals.user_id = $1
  AND goals.brand_id = $2
  AND goals.sale_type_id = $3
  AND goals.from_date = $4
  AND goals.to_date = $5
`

type GetUserBrandGoalParams struct {
	UserID     int32     `json:"user_id"`
	BrandID    int32     `json:"brand_id"`
	SaleTypeID int32     `json:"sale_type_id"`
	FromDate   time.Time `json:"from_date"`
	ToDate     time.Time `json:"to_date"`
}

func (q *Queries) GetUserBrandGoal(ctx context.Context, arg GetUserBrandGoalParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUserBrandGoal,
		arg.UserID,
		arg.BrandID,
		arg.SaleTypeID,
		arg.FromDate,
		arg.ToDate,
	)
	var value int64
	err := row.Scan(&value)
	return value, err
}

const getUserRank = `-- name: GetUserRank :one
WITH RankedUsers AS (SELECT user_id,
                            brand_id,
                            ratio,
                            ROW_NUMBER() OVER (ORDER BY ratio DESC) as rank
                     FROM user_brand_ratio
                     WHERE brand_id = $1
                       AND from_date = $2
                       AND to_date = $3)

SELECT rank
FROM RankedUsers
WHERE user_id = $4
`

type GetUserRankParams struct {
	BrandID  int32     `json:"brand_id"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	UserID   int32     `json:"user_id"`
}

func (q *Queries) GetUserRank(ctx context.Context, arg GetUserRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUserRank,
		arg.BrandID,
		arg.FromDate,
		arg.ToDate,
		arg.UserID,
	)
	var rank int64
	err := row.Scan(&rank)
	return rank, err
}

const insertUserBrandRatio = `-- name: InsertUserBrandRatio :exec
INSERT INTO user_brand_ratio (user_id, brand_id, ratio, from_date, to_date)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, brand_id, from_date, to_date)
DO
UPDATE SET ratio = EXCLUDED.ratio
`

type InsertUserBrandRatioParams struct {
	UserID   int32     `json:"user_id"`
	BrandID  int32     `json:"brand_id"`
	Ratio    float32   `json:"ratio"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

func (q *Queries) InsertUserBrandRatio(ctx context.Context, arg InsertUserBrandRatioParams) error {
	_, err := q.db.ExecContext(ctx, insertUserBrandRatio,
		arg.UserID,
		arg.BrandID,
		arg.Ratio,
		arg.FromDate,
		arg.ToDate,
	)
	return err
}
