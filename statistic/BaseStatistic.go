package statistic

import (
	sale "zhasa2.0/sale/entities"
)

type SaleSumByType map[sale.SaleType]sale.SaleAmount

func (s SaleSumByType) TotalSum() sale.SaleAmount {
	var totalSum sale.SaleAmount
	for _, amount := range s {
		totalSum += amount
	}
	return totalSum
}

type Statistic interface {
	ProvideSaleSums(period Period) SaleSumByType
}
