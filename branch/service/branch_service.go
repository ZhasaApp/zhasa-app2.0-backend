package service

import (
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/branch/entities"
	"zhasa2.0/branch/repository"
	. "zhasa2.0/manager/entities"
	. "zhasa2.0/sale/entities"
	. "zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
)

type BranchService interface {
	CreateBranch(request CreateBranchRequest) error
	GetBranches(period Period) ([]Branch, error)
	GetBranchYearStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error)
	GetBranchById(id BranchId) (*Branch, error)
	GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error)
	GetBranchGoal(period Period, id BranchId, typeId SaleTypeId) (SaleAmount, error)
	GetBranchRankedSalesManagers(period Period, branchId BranchId, pagination Pagination) (*[]SalesManager, error)
}

type DBBranchService struct {
	repo repository.BranchRepository
}

func (ds DBBranchService) GetBranchRankedSalesManagers(period Period, branchId BranchId, pagination Pagination) (*[]SalesManager, error) {
	from, to := period.ConvertToTime()
	return ds.repo.GetBranchRankedSalesManagers(from, to, branchId, pagination)
}

func (ds DBBranchService) GetBranchGoal(period Period, id BranchId, typeId SaleTypeId) (SaleAmount, error) {
	from, to := period.ConvertToTime()
	return ds.repo.GetBranchGoal(from, to, id, typeId)
}

func (ds DBBranchService) GetBranchSalesSums(from, to time.Time, branchId BranchId) (*SaleSumByType, error) {
	return ds.repo.GetBranchSalesSums(from, to, branchId)
}

func (ds DBBranchService) GetBranchYearStatistic(id BranchId, year int32) (*[]MonthlyYearStatistic, error) {
	return ds.repo.GetBranchYearMonthlyStatistic(id, year)
}

func (ds DBBranchService) CreateBranch(request CreateBranchRequest) error {
	return ds.repo.CreateBranch(request)
}

func (ds DBBranchService) GetBranches(period Period) ([]Branch, error) {
	return ds.repo.GetBranches(period)
}

func (ds DBBranchService) GetBranchById(id BranchId) (*Branch, error) {
	return ds.repo.GetBranchById(id)
}

func NewBranchService(repo repository.BranchRepository) BranchService {
	return DBBranchService{
		repo: repo,
	}
}
