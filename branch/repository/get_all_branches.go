package repository

import (
	"context"
	"database/sql"
	"fmt"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetAllBranches func() ([]entities.Branch, error)

func NewGetAllBranchesFunc(ctx context.Context, store generated.BranchStore) GetAllBranches {
	return func() ([]entities.Branch, error) {
		rows, err := store.GetAllBranches(ctx)
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
