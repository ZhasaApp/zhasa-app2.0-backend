package statistic

import (
	sale "zhasa2.0/sale/entities"
	"zhasa2.0/statistic/entities"
)

type SaleSumByType map[sale.SaleType]sale.SaleAmount

type Statistic interface {
	ProvideSaleSums(period entities.Period) SaleSumByType
}
