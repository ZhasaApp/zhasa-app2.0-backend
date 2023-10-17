package repository

import (
	"context"
	"database/sql"
	"fmt"
	"zhasa2.0/base"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
	"zhasa2.0/user/entities"
)

type GetUsersOrderedByRatioForGivenBrandFunc func(brandId int32, period statistic.Period, pagination base.Pagination) ([]entities.RatedUser, error)

func NewGetUsersOrderedByRatioForGivenBrandFunc(ctx context.Context, store generated.UserStore) GetUsersOrderedByRatioForGivenBrandFunc {
	return func(brandId int32, period statistic.Period, pagination base.Pagination) ([]entities.RatedUser, error) {
		from, to := period.ConvertToTime()
		params := generated.GetUsersOrderedByRatioForGivenBrandParams{
			BrandID:  brandId,
			FromDate: from,
			ToDate:   to,
			Limit:    pagination.PageSize,
			Offset:   pagination.GetOffset(),
			RoleID:   2,
		}
		rows, err := store.GetUsersOrderedByRatioForGivenBrand(ctx, params)
		ratedUsers := make([]entities.RatedUser, 0)
		if err == sql.ErrNoRows {
			fmt.Println(params)
			return ratedUsers, nil
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, row := range rows {
			ratedUsers = append(ratedUsers, entities.RatedUser{
				User: entities.User{
					Id:        row.ID,
					Phone:     "",
					Avatar:    row.AvatarUrl,
					FirstName: row.FirstName,
					LastName:  row.LastName,
				},
				Ratio: float64(row.Ratio) * 100,
				BranchInfo: entities.BranchInfo{
					Id:    row.BranchID,
					Title: row.BranchTitle,
				},
			})
		}
		fmt.Println(params)
		return ratedUsers, err
	}
}

type GetBranchUsersOrderedByRatioForGivenBrandFunc func(brandId int32, branchId int32, period statistic.Period, pagination base.Pagination) ([]entities.RatedUser, error)

func NewGetUsersOrderedBYRatioForGivenBrandAndBranchFunc(ctx context.Context, store generated.UserStore) GetBranchUsersOrderedByRatioForGivenBrandFunc {
	return func(brandId int32, branchId int32, period statistic.Period, pagination base.Pagination) ([]entities.RatedUser, error) {
		from, to := period.ConvertToTime()
		params := generated.GetBranchUsersOrderedByRatioForGivenBrandParams{
			BrandID:  brandId,
			FromDate: from,
			ToDate:   to,
			Offset:   pagination.GetOffset(),
			Limit:    pagination.PageSize,
			ID:       branchId,
			RoleID:   2,
		}
		rows, err := store.GetBranchUsersOrderedByRatioForGivenBrand(ctx, params)
		ratedUsers := make([]entities.RatedUser, 0)
		fmt.Println(params)
		if err == sql.ErrNoRows {
			fmt.Println(params)
			return ratedUsers, nil
		}
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, row := range rows {
			ratedUsers = append(ratedUsers, entities.RatedUser{
				User: entities.User{
					Id:        row.ID,
					Phone:     "",
					Avatar:    row.AvatarUrl,
					FirstName: row.FirstName,
					LastName:  row.LastName,
				},
				Ratio: float64(row.Ratio) * 100,
				BranchInfo: entities.BranchInfo{
					Id:    row.BranchID,
					Title: row.BranchTitle,
				},
			})
		}
		fmt.Println(params)
		return ratedUsers, err
	}
}
