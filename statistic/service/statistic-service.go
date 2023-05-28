package service

import (
	"log"
	. "zhasa2.0/base"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/statistic/entities"
	. "zhasa2.0/statistic/repository"
)

type AnalyticsService interface {
	ProvideRankedManagers(pagination Pagination, period Period) (*[]SalesManager, error)
}

type DBAnalyticsService struct {
	rankingsRepo RankingsRepository
}

func (D DBAnalyticsService) ProvideRankedManagers(pagination Pagination, period Period) (*[]SalesManager, error) {
	from, to := period.ConvertToTime()
	log.Println(from)
	log.Println(to)
	return D.rankingsRepo.ProvideRankedManagers(pagination, from, to)
}

func NewAnalyticsService(rankingsRepo RankingsRepository) AnalyticsService {
	return DBAnalyticsService{
		rankingsRepo: rankingsRepo,
	}
}
