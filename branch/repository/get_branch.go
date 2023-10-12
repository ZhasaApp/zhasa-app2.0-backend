package repository

import (
	"context"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchByIdFunc func(id int32) (*entities.BranchInfo, error)

func NewGetBranchByIdFunc(ctx context.Context, store generated.BranchStore) GetBranchByIdFunc {
	return func(id int32) (*entities.BranchInfo, error) {
		row, err := store.GetBranchById(ctx, id)
		if err != nil {
			return nil, err
		}
		return &entities.BranchInfo{
			BranchId:    row.ID,
			Title:       row.Title,
			Description: row.Description,
		}, nil
	}
}
