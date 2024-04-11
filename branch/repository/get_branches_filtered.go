package repository

import (
	"context"
	"database/sql"
	"fmt"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchesFiltered func(params generated.GetBranchesSearchParams) ([]entities.Branch, error)

func NewGetBranchesFiltered(ctx context.Context, store generated.BranchStore) GetBranchesFiltered {
	return func(params generated.GetBranchesSearchParams) ([]entities.Branch, error) {
		rows, err := store.GetBranchesSearch(ctx, params)
		if err != nil {
			return nil, err
		}

		branches := make([]entities.Branch, 0)
		if err == sql.ErrNoRows {
			return branches, nil
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, row := range rows {
			branches = append(branches, entities.Branch{
				BranchId:    row.ID,
				Title:       row.Title,
				Description: row.Description,
			})
		}
		return branches, nil
	}
}
