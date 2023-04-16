package statistic

import (
	entities2 "zhasa2.0/manager/entities"
	sale "zhasa2.0/sale/entities"
	"zhasa2.0/statistic/entities"
)

type SaleSumByType map[sale.SaleType]sale.SaleAmount

type Statistic interface {
	Calculate() SaleSumByType
}

type SalesManagerStatistic struct {
	p  entities.Period
	sm entities2.SalesManager
}

type BranchStatistic struct {
	p entities.Period
}
