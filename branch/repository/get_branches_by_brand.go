package repository

import (
	"context"
	"errors"
	"zhasa2.0/branch/entities"
	generated "zhasa2.0/db/sqlc"
)

type GetBranchesByBrandFunc func(brandId int32) ([]entities.BranchDescriptionInfo, error)

func NewGetBranchesByBrandFunc(ctx context.Context, store generated.BranchStore) GetBranchesByBrandFunc {
	return func(brandId int32) ([]entities.BranchDescriptionInfo, error) {
		branches, err := store.GetBranchesByBrandId(ctx, brandId)
		if err != nil {
			return nil, errors.New("branches not found for given brand")
		}
		var result []entities.BranchDescriptionInfo
		for _, branch := range branches {
			result = append(result, entities.BranchDescriptionInfo{
				BranchId:    branch.ID,
				Title:       branch.Title,
				Description: branch.Description,
			})
		}
		return result, nil
	}
}
