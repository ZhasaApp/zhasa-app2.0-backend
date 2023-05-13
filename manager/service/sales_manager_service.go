package service

import (
	"time"
	. "zhasa2.0/manager/entities"
	"zhasa2.0/manager/repository"
	sale "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic"
	. "zhasa2.0/statistic/entities"
)

type SalesManagerService interface {
	GetSalesManagerByUserId(userId int32) (*SalesManager, error)
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(sale sale.Sale) error
	GetSalesManagerGoal(from, to time.Time, salesManagerId SalesManagerId) (sale.SaleAmount, error)
	GetSalesManagerSums(from, to time.Time, salesManagerId SalesManagerId) (*SaleSumByType, error)
	GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]YearStatisticByMonth, error)
}

type DBSalesManagerService struct {
	repo repository.SalesManagerRepository
	repository2.SaleTypeRepository
}

func (dbs DBSalesManagerService) SaveSale(sale sale.Sale) error {
	return dbs.repo.SaveSale(sale.SaleManagerId, sale.SaleDate, sale.SalesAmount, sale.SalesTypeId)
}

func NewSalesManagerService(repo repository.SalesManagerRepository, saleTypeRepo repository2.SaleTypeRepository) SalesManagerService {
	return DBSalesManagerService{
		repo:               repo,
		SaleTypeRepository: saleTypeRepo,
	}
}

type salesManagerVisitor interface {
	provideManagers() (error, *SalesManagers)
}

func (dbs DBSalesManagerService) ProvideManagers(visitor salesManagerVisitor) (error, *SalesManagers) {
	return visitor.provideManagers()
}

func (dbs DBSalesManagerService) CreateSalesManager(userId int32, branchId int32) error {
	return dbs.repo.CreateSalesManager(userId, branchId)
}

func (dbs DBSalesManagerService) GetSalesManagerByUserId(userId int32) (*SalesManager, error) {
	return dbs.repo.GetSalesManagerByUserId(userId)
}

func (dbs DBSalesManagerService) GetSalesManagerGoal(fromDate time.Time, to time.Time, salesManagerId SalesManagerId) (sale.SaleAmount, error) {
	return dbs.repo.GetSalesManagerGoalAmount(salesManagerId, fromDate, to)
}

func (dbs DBSalesManagerService) GetSalesManagerSums(from, to time.Time, salesManagerId SalesManagerId) (*SaleSumByType, error) {
	return dbs.repo.ProvideSums(salesManagerId, from, to)
}

func (dbs DBSalesManagerService) GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]YearStatisticByMonth, error) {
	return dbs.repo.GetMonthlyYearSaleStatistic(smId, year)
}
