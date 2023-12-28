package repository

import (
	. "zhasa2.0/sale/repository"
	"zhasa2.0/statistic"
	"zhasa2.0/statistic/entities"
)

type GetBrandMonthlyYearStatisticFunc func(year int32, brandId int32, saleTypeID int32) ([]entities.MonthlyYearStatistic, error)

func NewGetBrandMonthlyYearStatisticFunc(
	saleTypeRepo SaleTypeRepository,
	brandOverallGoalFunc GetBrandOverallGoalFunc,
	brandSaleSumFunc GetBrandSaleSumFunc,
) GetBrandMonthlyYearStatisticFunc {
	return func(year int32, brandId int32, saleTypeID int32) ([]entities.MonthlyYearStatistic, error) {
		saleType, err := saleTypeRepo.GetSaleType(saleTypeID)
		if err != nil {
			return nil, err
		}
		result := make([]entities.MonthlyYearStatistic, 0)

		for month := 1; month <= 12; month++ {
			period := statistic.MonthPeriod{
				MonthNumber: int32(month),
				Year:        year,
			}
			goal, err := brandOverallGoalFunc(brandId, saleType.Id, period)

			sum, err := brandSaleSumFunc(brandId, saleType.Id, period)
			if err != nil {
				return nil, err
			}

			stat := entities.MonthlyYearStatistic{
				SaleType: *saleType,
				Month:    int32(month),
				Amount:   sum,
				Goal:     goal,
			}
			result = append(result, stat)
		}

		return result, nil
	}
}
