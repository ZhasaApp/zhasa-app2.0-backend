package entities

import (
	. "zhasa2.0/sale/entities"
)

type MonthlyYearStatistic struct {
	SaleType SaleType
	Month    MonthNumber
	Amount   SaleAmount
}
