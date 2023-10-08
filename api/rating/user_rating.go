package rating

import (
	"context"
	"fmt"
	generated "zhasa2.0/db/sqlc"
	"zhasa2.0/statistic/entities"
)

type GetUserRatingFunc func(userId int32, brandId int32, period entities.Period) (int32, error)

func NewGetUserRatingFunc(ctx context.Context, store generated.UserBrandStore) GetUserRatingFunc {
	return func(userId int32, brandId int32, period entities.Period) (int32, error) {
		from, to := period.ConvertToTime()
		rating, err := store.GetUserRank(ctx, generated.GetUserRankParams{
			BrandID:  brandId,
			FromDate: from,
			ToDate:   to,
			UserID:   userId,
		})
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		return int32(rating), nil
	}
}
