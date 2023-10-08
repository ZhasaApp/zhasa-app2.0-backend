package repository

import (
	"errors"
	"fmt"
	"zhasa2.0/api/rating"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/sale/repository"
	"zhasa2.0/statistic/entities"
)

type CalculateUserBrandRatio func(userId int32, brandId int32, period entities.Period) (float32, error)

func NewCalculateUserBrandRatio(saleTypeRepo repository.SaleTypeRepository, saleRepo repository.SaleRepository, goalFunc UserBrandGoalFunc, brandFunc GetUserBrandFunc) CalculateUserBrandRatio {
	return func(userId int32, brandId int32, period entities.Period) (float32, error) {
		var goalAchievementPercent float32
		ratioRows := make([]rating.RatioRow, 0)
		userBrand, err := brandFunc(userId, brandId)

		if err != nil {
			fmt.Println(errors.New("user brand not found"))
			return 0, err
		}

		from, to := period.ConvertToTime()

		saleTypes, err := saleTypeRepo.GetSaleTypes()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		for _, saleType := range *saleTypes {
			amount, err := saleRepo.GetSumByUserIdBrandIdPeriodSaleTypeId(generated.GetSaleSumByUserIdBrandIdPeriodSaleTypeIdParams{
				ID:         userId,
				BrandID:    brandId,
				SaleDate:   from,
				SaleDate_2: to,
				SaleTypeID: saleType.Id,
			})
			if err != nil {
				fmt.Println(err)
				return 0, err
			}
			goal := goalFunc(generated.GetUserBrandGoalParams{
				UserBrand:  userBrand,
				SaleTypeID: saleType.Id,
				FromDate:   from,
				FromDate_2: to,
			})

			ratioRows = append(ratioRows, rating.RatioRow{
				Amount:  amount,
				Goal:    goal,
				Gravity: saleType.Gravity,
			})
		}

		goalAchievementPercent = rating.CalculateRatio(ratioRows)

		return goalAchievementPercent, nil
	}
}
