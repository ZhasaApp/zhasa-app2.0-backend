package repository

import (
	"time"
	"zhasa2.0/manager/entities"
	"zhasa2.0/statistic"
)

type StatisticRepository interface {
	ProvideSums(salesManagerId entities.SalesManagerId, from time.Time, to time.Time) (*statistic.SaleSumByType, error)
}
