package repository

import (
	"time"
	"zhasa2.0/statistic"
)

type StatisticRepository interface {
	ProvideSums(from time.Time, to time.Time) (*statistic.SaleSumByType, error)
}
