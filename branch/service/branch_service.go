package service

import (
	. "zhasa2.0/branch/entities"
	"zhasa2.0/branch/repository"
	. "zhasa2.0/statistic/entities"
)

type BranchService interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranches() ([]Branch, error)
	GetBranchYearStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error)
	GetBranchById(id BranchId) (Branch, error)
}

type DBBranchService struct {
	repo repository.BranchRepository
}

func (ds DBBranchService) GetBranchYearStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error) {
	return ds.repo.GetBranchYearMonthlyStatistic(id, year)
}

func (ds DBBranchService) CreateBranch(request CreateBranchRequest) error {
	return ds.repo.CreateBranch(request)
}

func (ds DBBranchService) GetBranches() ([]Branch, error) {
	return ds.repo.GetBranches()
}

func (ds DBBranchService) GetBranchById(id BranchId) (Branch, error) {
	return ds.GetBranchById(id)
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return DBBranchService{
		repo: repo,
	}
}
