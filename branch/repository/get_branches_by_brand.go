package repository

import (
	"context"
	"errors"
	"sort"
	"zhasa2.0/branch/entities"
	"zhasa2.0/brand"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/rating"
	entities2 "zhasa2.0/sale/entities"
	"zhasa2.0/sale/repository"
	"zhasa2.0/statistic"
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

type RatedBranchesFunc func(brandId int32, period statistic.Period, saleTypes []entities2.SaleType) ([]entities.Branch, error)

func NewRatedBranchesFunc(
	ctx context.Context,
	store generated.BranchStore,
	getBranchBrandSaleSumFunc repository.GetBranchBrandSaleSumFunc,
	getBranchBrandGoalFunc brand.GetBranchBrandGoalFunc,
) RatedBranchesFunc {
	return func(brandId int32, period statistic.Period, saleTypes []entities2.SaleType) ([]entities.Branch, error) {
		branches, err := store.GetBranchesByBrandId(ctx, brandId)
		if err != nil {
			return nil, errors.New("branches not found for given brand")
		}
		var result []entities.Branch
		for _, branch := range branches {
			branchRatioRows := make([]rating.RatioRow, 0)
			for _, saleType := range saleTypes {
				salesSum, _ := getBranchBrandSaleSumFunc(branch.ID, brandId, saleType.Id, period)

				goal, _ := getBranchBrandGoalFunc(branch.ID, brandId, saleType.Id, period)

				branchRatioRows = append(branchRatioRows, rating.RatioRow{
					Achieved: salesSum,
					Goal:     goal,
					Gravity:  saleType.Gravity,
				})
			}
			result = append(result, entities.Branch{
				BranchId:        branch.ID,
				Title:           branch.Title,
				Description:     branch.Description,
				GoalAchievement: rating.CalculateRatio(branchRatioRows) * 100,
			})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].GoalAchievement > result[j].GoalAchievement
		})
		return result, nil
	}
}
