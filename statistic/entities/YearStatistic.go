package entities

import (
	. "zhasa2.0/sale/entities"
)

type YearStatisticByMonth struct {
	SaleType SaleType
	Month    MonthNumber
	Amount   SaleAmount
}
