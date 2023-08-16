// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: branch_director.sql

package generated

import (
	"context"
	"time"
)

const createBranchDirector = `-- name: CreateBranchDirector :one
INSERT INTO branch_directors (user_id, branch_id)
VALUES ($1, $2) RETURNING id
`

type CreateBranchDirectorParams struct {
	UserID   int32 `json:"user_id"`
	BranchID int32 `json:"branch_id"`
}

func (q *Queries) CreateBranchDirector(ctx context.Context, arg CreateBranchDirectorParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, createBranchDirector, arg.UserID, arg.BranchID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createSalesManagerGoalByType = `-- name: CreateSalesManagerGoalByType :exec
INSERT INTO sales_manager_goals_by_types (sales_manager_id, from_date, to_date, amount, type_id)
VALUES ($1, $2, $3, $4, $5)
`

type CreateSalesManagerGoalByTypeParams struct {
	SalesManagerID int32     `json:"sales_manager_id"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Amount         int64     `json:"amount"`
	TypeID         int32     `json:"type_id"`
}

func (q *Queries) CreateSalesManagerGoalByType(ctx context.Context, arg CreateSalesManagerGoalByTypeParams) error {
	_, err := q.db.ExecContext(ctx, createSalesManagerGoalByType,
		arg.SalesManagerID,
		arg.FromDate,
		arg.ToDate,
		arg.Amount,
		arg.TypeID,
	)
	return err
}

const getBranchDirectorByBranchId = `-- name: GetBranchDirectorByBranchId :one
SELECT user_id, phone, first_name, last_name, avatar_url, branch_director_id, branch_id, branch_title
FROM branch_directors_view bdv
WHERE bdv.branch_id = $1
`

func (q *Queries) GetBranchDirectorByBranchId(ctx context.Context, branchID int32) (BranchDirectorsView, error) {
	row := q.db.QueryRowContext(ctx, getBranchDirectorByBranchId, branchID)
	var i BranchDirectorsView
	err := row.Scan(
		&i.UserID,
		&i.Phone,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.BranchDirectorID,
		&i.BranchID,
		&i.BranchTitle,
	)
	return i, err
}

const getBranchDirectorByUserId = `-- name: GetBranchDirectorByUserId :many
SELECT user_id, phone, first_name, last_name, avatar_url, branch_director_id, branch_id, branch_title
FROM branch_directors_view bdv
WHERE bdv.user_id = $1
`

func (q *Queries) GetBranchDirectorByUserId(ctx context.Context, userID int32) ([]BranchDirectorsView, error) {
	rows, err := q.db.QueryContext(ctx, getBranchDirectorByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BranchDirectorsView
	for rows.Next() {
		var i BranchDirectorsView
		if err := rows.Scan(
			&i.UserID,
			&i.Phone,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.BranchDirectorID,
			&i.BranchID,
			&i.BranchTitle,
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

const getBranchGoalBySaleType = `-- name: GetBranchGoalBySaleType :one
SELECT COALESCE(amount, 0)
FROM branch_goals_by_types
WHERE from_date = $1
  AND to_date = $2
  AND branch_id = $3
  AND type_id = $4
`

type GetBranchGoalBySaleTypeParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	BranchID int32     `json:"branch_id"`
	TypeID   int32     `json:"type_id"`
}

func (q *Queries) GetBranchGoalBySaleType(ctx context.Context, arg GetBranchGoalBySaleTypeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getBranchGoalBySaleType,
		arg.FromDate,
		arg.ToDate,
		arg.BranchID,
		arg.TypeID,
	)
	var amount int64
	err := row.Scan(&amount)
	return amount, err
}

const getSMGoal = `-- name: GetSMGoal :one
SELECT COALESCE(amount, 0)
FROM sales_manager_goals_by_types
WHERE from_date = $1
  AND to_date = $2
  AND sales_manager_id = $3
  AND type_id = $4
`

type GetSMGoalParams struct {
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	SalesManagerID int32     `json:"sales_manager_id"`
	TypeID         int32     `json:"type_id"`
}

func (q *Queries) GetSMGoal(ctx context.Context, arg GetSMGoalParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSMGoal,
		arg.FromDate,
		arg.ToDate,
		arg.SalesManagerID,
		arg.TypeID,
	)
	var amount int64
	err := row.Scan(&amount)
	return amount, err
}

const setBranchGoalBySaleType = `-- name: SetBranchGoalBySaleType :exec
INSERT INTO branch_goals_by_types (from_date, to_date, amount, branch_id, type_id)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (from_date, to_date, sales_manager_id, type_id)
DO
UPDATE SET amount = EXCLUDED.amount
`

type SetBranchGoalBySaleTypeParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Amount   int64     `json:"amount"`
	BranchID int32     `json:"branch_id"`
	TypeID   int32     `json:"type_id"`
}

func (q *Queries) SetBranchGoalBySaleType(ctx context.Context, arg SetBranchGoalBySaleTypeParams) error {
	_, err := q.db.ExecContext(ctx, setBranchGoalBySaleType,
		arg.FromDate,
		arg.ToDate,
		arg.Amount,
		arg.BranchID,
		arg.TypeID,
	)
	return err
}

const setSmGoalBySaleType = `-- name: SetSmGoalBySaleType :exec
INSERT INTO sales_manager_goals_by_types (from_date, to_date, amount, sales_manager_id, type_id)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (from_date, to_date, sales_manager_id, type_id)
DO
UPDATE SET amount = EXCLUDED.amount
`

type SetSmGoalBySaleTypeParams struct {
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Amount         int64     `json:"amount"`
	SalesManagerID int32     `json:"sales_manager_id"`
	TypeID         int32     `json:"type_id"`
}

func (q *Queries) SetSmGoalBySaleType(ctx context.Context, arg SetSmGoalBySaleTypeParams) error {
	_, err := q.db.ExecContext(ctx, setSmGoalBySaleType,
		arg.FromDate,
		arg.ToDate,
		arg.Amount,
		arg.SalesManagerID,
		arg.TypeID,
	)
	return err
}
