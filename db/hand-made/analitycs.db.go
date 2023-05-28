package hand_made

import (
	"context"
	"time"
)

const getRankedSalesManagers = `-- name: GetRankedSalesManagers :many
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
ORDER BY rating_position LIMIT $3
OFFSET $4;
`

type GetRankedSalesManagersParams struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

// GetRankedSalesManagers get the ranked sales managers by their total sales divided by their sales goal amount for the given period.
func (d DBCustomQuerier) GetRankedSalesManagers(ctx context.Context, arg GetRankedSalesManagersParams) ([]GetRankedSalesManagersRow, error) {
	rows, err := d.db.QueryContext(ctx, getRankedSalesManagers,
		arg.FromDate,
		arg.ToDate,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRankedSalesManagersRow
	for rows.Next() {
		var i GetRankedSalesManagersRow
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
