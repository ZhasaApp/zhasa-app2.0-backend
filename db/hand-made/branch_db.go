package hand_made

import (
	"context"
	"time"
)

const getBranchYearStatistic = `-- name: GetBranchYearStatistic :many
SELECT
    st.id AS sale_type_id,
  CAST(EXTRACT(MONTH FROM s.sale_date) AS INTEGER) AS month_number,
    SUM(s.amount) AS total_amount
FROM
    sales AS s
        JOIN
    sale_types AS st ON s.sale_type_id = st.id
        JOIN
    sales_managers AS sm ON s.sales_manager_id = sm.id
WHERE
        sm.branch_id = $1 AND
        DATE_PART('year', s.sale_date)::integer = $2
GROUP BY
    st.id,
   	month_number 
ORDER BY
    month_number,
    st.title
`

type GetBranchYearStatisticParams struct {
	BranchID   int32 `json:"branch_id"`
	YearNumber int32 `json:"year_number"`
}

type GetBranchYearStatisticRow struct {
	SaleTypeId  int32 `json:"sale_type_id"`
	MonthNumber int32 `json:"month_number"`
	TotalAmount int64 `json:"total_amount"`
}

func (d DBCustomQuerier) GetBranchYearStatistic(ctx context.Context, arg GetBranchYearStatisticParams) ([]GetBranchYearStatisticRow, error) {
	rows, err := d.db.QueryContext(ctx, getBranchYearStatistic, arg.BranchID, arg.YearNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBranchYearStatisticRow
	for rows.Next() {
		var i GetBranchYearStatisticRow
		if err := rows.Scan(&i.SaleTypeId, &i.MonthNumber, &i.TotalAmount); err != nil {
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

const getBranchRankedSalesManagers = `-- name: GetBranchRankedSalesManagers :many
WITH goal_sales AS (SELECT 
                        sm.user_id AS user_id,
                           sm.sales_manager_id        AS sales_manager_id,
                           sm.first_name              AS first_name,
                           sm.last_name               AS last_name,
                           COALESCE(smg.amount ,0)                AS sale_goal,
                           sm.branch_id               AS branch_id,
                           COALESCE(SUM(s.amount), 0) AS total_sales_sum
                    FROM sales_managers_view sm
                             LEFT JOIN
                         sales_manager_goals smg
                         ON sm.sales_manager_id = smg.sales_manager_id
                             AND smg.from_date = $1
                             AND smg.to_date = $2
                             LEFT JOIN
                         sales s
                         ON sm.sales_manager_id = s.sales_manager_id
                             AND s.sale_date BETWEEN $1 AND $2
                    GROUP BY  sm.user_id,
        						sm.sales_manager_id,
        sm.first_name,
        sm.last_name,
        sm.branch_id,
        smg.amount),
     rankings AS (SELECT user_id, sales_manager_id, first_name, last_name, sale_goal, branch_id, total_sales_sum,
                         CASE
                            WHEN sale_goal IS NULL OR sale_goal = 0 THEN 0
                             ELSE total_sales_sum::decimal / sale_goal:: decimal
END
AS ratio,
        RANK() OVER (ORDER BY CASE
                             WHEN sale_goal IS NULL OR sale_goal = 0 THEN 0
                             ELSE total_sales_sum::decimal / sale_goal::decimal
                         END DESC) AS rating_position
    FROM
        goal_sales
)
SELECT user_id,
       sales_manager_id,
       first_name,
       last_name,
       branch_id,
       sale_goal,
       total_sales_sum,
       ratio,
       rating_position
FROM rankings
WHERE branch_id = $5
ORDER BY rating_position LIMIT $3
OFFSET $4;
`

type GetBranchRankedSalesManagersParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
	BranchID int32     `json:"branch_id"`
}

type GetBranchRankedSalesManagersRow struct {
	UserId         int32   `json:"user_id"`
	SalesManagerID int32   `json:"sales_manager_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	BranchID       int32   `json:"branch_id"`
	SaleGoal       int64   `json:"sale_goal"`
	TotalSalesSum  int64   `json:"total_sales_sum"`
	Ratio          float64 `json:"ratio"`
	RatingPosition int32   `json:"rating_position"`
}

// GetBranchRankedSalesManagers get the ranked sales managers by their total sales divided by their sales goal amount for the given period.
func (d DBCustomQuerier) GetBranchRankedSalesManagers(ctx context.Context, arg GetBranchRankedSalesManagersParams) ([]GetBranchRankedSalesManagersRow, error) {
	rows, err := d.db.QueryContext(ctx, getBranchRankedSalesManagers,
		arg.FromDate,
		arg.ToDate,
		arg.Limit,
		arg.Offset,
		arg.BranchID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBranchRankedSalesManagersRow
	for rows.Next() {
		var i GetBranchRankedSalesManagersRow
		if err := rows.Scan(
			&i.UserId,
			&i.SalesManagerID,
			&i.FirstName,
			&i.LastName,
			&i.BranchID,
			&i.SaleGoal,
			&i.TotalSalesSum,
			&i.Ratio,
			&i.RatingPosition,
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
