package repository

import (
	"fmt"
	"zhasa2.0/rating"
	"zhasa2.0/statistic"
	"zhasa2.0/user/repository"
)

type CalculateUserBrandRatio func(userId int32, brandId int32, period statistic.Period) (float32, error)

func NewCalculateUserBrandRatio(saleTypeRepo SaleTypeRepository, userSaleSum GetSaleSumByUserBrandTypePeriodFunc, goalFunc repository.UserBrandGoalFunc) CalculateUserBrandRatio {
	return func(userId int32, brandId int32, period statistic.Period) (float32, error) {
		var goalAchievementPercent float32
		ratioRows := make([]rating.RatioRow, 0)

		saleTypes, err := saleTypeRepo.GetSaleTypes()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		for _, saleType := range *saleTypes {
			amount, err := userSaleSum(userId, brandId, saleType.Id, period)
			if err != nil {
				fmt.Println(err)
				return 0, err
			}
			goal := goalFunc(userId, brandId, saleType.Id, period)

			ratioRows = append(ratioRows, rating.RatioRow{
				Achieved: amount,
				Goal:     goal,
				Gravity:  saleType.Gravity,
			})
		}

		goalAchievementPercent = rating.CalculateRatio(ratioRows)

		return goalAchievementPercent, nil
	}
}
