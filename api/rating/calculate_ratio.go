package rating

type RatioRow struct {
	Amount  int64
	Goal    int64
	Gravity int32
}

func CalculateRatio(rows []RatioRow) float32 {
	var totalWeightedRatio float32
	var totalGravity int32

	for _, row := range rows {
		var ratio float32
		if row.Goal != 0 {
			ratio = float32(row.Amount) / float32(row.Goal)
			totalGravity += row.Gravity
		}
		weightedRatio := ratio * float32(row.Gravity)
		totalWeightedRatio += weightedRatio
	}
	if totalGravity == 0 {
		return 0
	}
	return totalWeightedRatio / float32(totalGravity)
}
