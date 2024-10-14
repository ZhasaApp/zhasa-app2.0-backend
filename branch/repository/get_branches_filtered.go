package repository

import (
	"context"
	"database/sql"
	"fmt"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchesFilteredAsc func(params generated.GetBranchesSearchAscParams) ([]entities.Branch, error)

func NewGetBranchesFilteredAsc(ctx context.Context, store generated.BranchStore) GetBranchesFilteredAsc {
	return func(params generated.GetBranchesSearchAscParams) ([]entities.Branch, error) {
		params.Offset = (params.Offset - 1) * params.Limit
		rows, err := store.GetBranchesSearchAsc(ctx, params)
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

type GetBranchesFilteredDesc func(params generated.GetBranchesSearchDescParams) ([]entities.Branch, error)

func NewGetBranchesFilteredDesc(ctx context.Context, store generated.BranchStore) GetBranchesFilteredDesc {
	return func(params generated.GetBranchesSearchDescParams) ([]entities.Branch, error) {
		params.Offset = (params.Offset - 1) * params.Limit
		rows, err := store.GetBranchesSearchDesc(ctx, params)
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

type GetBranchesFilteredCount func(search string) (int64, error)

func NewGetBranchesFilteredCount(ctx context.Context, store generated.BranchStore) GetBranchesFilteredCount {
	return func(search string) (int64, error) {
		count, err := store.GetBranchesSearchCount(ctx, search)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}
