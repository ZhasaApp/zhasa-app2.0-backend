package entities

import (
	. "zhasa2.0/sale/entities"
)

type MonthlyYearStatistic struct {
	SaleType SaleType
	Month    int32
	Amount   int64
	Goal     int64
}
