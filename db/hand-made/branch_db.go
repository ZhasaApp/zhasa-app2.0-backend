package hand_made

import (
	"context"
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
