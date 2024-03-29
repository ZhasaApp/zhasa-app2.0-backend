package service

//
//import (
//	. "zhasa2.0/base"
//	"zhasa2.0/manager/entities"
//	. "zhasa2.0/sale/entities"
//	. "zhasa2.0/statistic/entities"
//)
//
//type RatioRow struct {
//	Achieved  SaleAmount
//	Goal    SaleAmount
//	Gravity int32
//}
//
//func (dbs DBSalesManagerService) UpdateRatio(smId entities.SalesManagerId, period Period) (Percent, error) {
//	types, err := dbs.GetSaleTypes()
//	from, to := period.ConvertToTime()
//	if err != nil {
//		return 0, err
//	}
//
//	ratioRows := make([]RatioRow, 0)
//
//	for _, item := range *types {
//		sum, _ := dbs.statisticRepo.GetSalesSumBySaleTypeAndManager(smId, item.Id, from, to)
//		goal, _ := dbs.statisticRepo.GetSalesGoalBySaleTypeAndManager(smId, item.Id, from, to)
//
//		ratioRows = append(ratioRows, RatioRow{
//			Achieved:  sum,
//			Goal:    goal,
//			Gravity: item.Gravity,
//		})
//	}
//
//	ratio := Percent(CalculateRatio(ratioRows))
//
//	err = dbs.statisticRepo.SetRatioByPeriod(smId, ratio, from, to)
//
//	return Percent(CalculateRatio(ratioRows)), err
//}
//
//func CalculateRatio(rows []RatioRow) float32 {
//	var totalWeightedRatio float32
//	var totalGravity int32
//
//	for _, row := range rows {
//		var ratio float32
//		if row.Goal != 0 {
//			ratio = float32(row.Achieved) / float32(row.Goal)
//			totalGravity += row.Gravity
//		}
//		weightedRatio := ratio * float32(row.Gravity)
//		totalWeightedRatio += weightedRatio
//	}
//	if totalGravity == 0 {
//		return 0
//	}
//	return totalWeightedRatio / float32(totalGravity)
//}
