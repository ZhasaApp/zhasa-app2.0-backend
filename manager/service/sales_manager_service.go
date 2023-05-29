package service

import (
	"time"
	. "zhasa2.0/base"
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
	SaveSale(sale sale.Sale) (*sale.Sale, error)
	GetSalesManagerGoal(from, to time.Time, salesManagerId SalesManagerId) (sale.SaleAmount, error)
	GetSalesManagerSums(from, to time.Time, salesManagerId SalesManagerId) (*SaleSumByType, error)
	GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error)
	GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]sale.Sale, error)
	GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, period Period) (*[]sale.Sale, error)
	GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error)
}

type DBSalesManagerService struct {
	repo repository.SalesManagerRepository
	repository2.SaleTypeRepository
}

func (dbs DBSalesManagerService) GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, period Period) (*[]sale.Sale, error) {
	from, to := period.ConvertToTime()
	return dbs.repo.GetManagerSalesByPeriod(salesManagerId, pagination, from, to)
}

func (dbs DBSalesManagerService) GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error) {
	return dbs.repo.GetSalesManagerSalesCount(salesManagerId)
}

func (dbs DBSalesManagerService) GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]sale.Sale, error) {
	return dbs.repo.GetManagerSales(salesManagerId, pagination)
}

func (dbs DBSalesManagerService) SaveSale(sale sale.Sale) (*sale.Sale, error) {
	return dbs.repo.SaveSale(sale.SaleManagerId, sale.SaleDate, sale.SalesAmount, sale.SaleType.Id, sale.SaleDescription)
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

func (dbs DBSalesManagerService) GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error) {
	return dbs.repo.GetMonthlyYearSaleStatistic(smId, year)
}
