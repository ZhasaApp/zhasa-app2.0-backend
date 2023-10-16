package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"zhasa2.0/brand"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic"
	"zhasa2.0/statistic/entities"
)

type GetBranchBrandMonthlyYearStatisticFunc func(year int32, branchId int32, brandId int32) ([]entities.MonthlyYearStatistic, error)

func NewGetBranchBrandMonthlyYearStatisticFunc(ctx context.Context, saleTypeRepo SaleTypeRepository, branchBrandGoalFunc brand.GetBranchBrandGoalFunc, branchBrandFunc brand.GetBranchBrandFunc, store generated.SaleStore) GetBranchBrandMonthlyYearStatisticFunc {
	return func(year int32, branchId int32, brandId int32) ([]entities.MonthlyYearStatistic, error) {
		saleTypes, err := saleTypeRepo.GetSaleTypes()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result := make([]entities.MonthlyYearStatistic, 0)

		branchBrand, err := branchBrandFunc(branchId, brandId)
		if err != nil {
			log.Println(err)
			return nil, errors.New("branch brand not found")
		}

		for _, saleType := range *saleTypes {
			for month := 1; month <= 12; month++ {
				period := statistic.MonthPeriod{
					MonthNumber: int32(month),
					Year:        year,
				}
				from, to := period.ConvertToTime()
				goal, err := branchBrandGoalFunc(branchBrand, saleType.Id, period)

				sum, err := store.GetBranchBrandSaleSumByGivenDateRange(ctx, generated.GetBranchBrandSaleSumByGivenDateRangeParams{
					BranchID:   branchId,
					BrandID:    brandId,
					SaleTypeID: saleType.Id,
					SaleDate:   from,
					SaleDate_2: to,
				})
				if err != nil {
					log.Println(err)
				}

				stat := entities.MonthlyYearStatistic{
					SaleType: saleType,
					Month:    int32(month),
					Amount:   GetInt64FromInterface(sum),
					Goal:     goal,
				}
				result = append(result, stat)
			}

		}
		return result, nil
	}
}
