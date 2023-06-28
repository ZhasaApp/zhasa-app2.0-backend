package service

import (
	"errors"
	. "zhasa2.0/base"
	"zhasa2.0/manager/entities"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/statistic/entities"
)

type RatioRow struct {
	amount  SaleAmount
	goal    SaleAmount
	gravity int32
}

func (dbs DBSalesManagerService) UpdateRatio(smId entities.SalesManagerId, period Period) (Percent, error) {
	types, err := dbs.GetSaleTypes()
	from, to := period.ConvertToTime()
	if err != nil {
		return 0, err
	}

	ratioRows := make([]RatioRow, 0)

	for _, item := range *types {
		sum, _ := dbs.statisticRepo.GetSalesSumBySaleTypeAndManager(smId, item.Id, from, to)
		goal, err := dbs.statisticRepo.GetSalesGoalBySaleTypeAndManager(smId, item.Id, from, to)
		if err != nil || goal != 0 {
			return 0, errors.New("no goal found for given sale type: " + string(item.Id))
		}

		ratioRows = append(ratioRows, RatioRow{
			amount:  sum,
			goal:    goal,
			gravity: item.Gravity,
		})
	}

	ratio := Percent(CalculateRatio(ratioRows))

	err = dbs.statisticRepo.SetRatioByPeriod(smId, ratio, from, to)

	return Percent(CalculateRatio(ratioRows)), err
}

func CalculateRatio(rows []RatioRow) float32 {
	var totalWeightedRatio float32
	var totalGravity int32

	for _, row := range rows {
		ratio := float32(row.amount) / float32(row.goal)
		weightedRatio := ratio * float32(row.gravity)
		totalWeightedRatio += weightedRatio
		totalGravity += row.gravity
	}

	return totalWeightedRatio / float32(totalGravity)
}
