package repository

import (
	"context"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchByIdFunc func(id int32) (*entities.BranchDescriptionInfo, error)

func NewGetBranchByIdFunc(ctx context.Context, store generated.BranchStore) GetBranchByIdFunc {
	return func(id int32) (*entities.BranchDescriptionInfo, error) {
		row, err := store.GetBranchById(ctx, id)
		if err != nil {
			return nil, err
		}
		return &entities.BranchDescriptionInfo{
			BranchId:    row.ID,
			Title:       row.Title,
			Description: row.Description,
		}, nil
	}
}
