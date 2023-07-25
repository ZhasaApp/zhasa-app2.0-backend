package service

import (
	"log"
	"time"
	. "zhasa2.0/base"
	. "zhasa2.0/manager/entities"
	"zhasa2.0/manager/repository"
	. "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
	. "zhasa2.0/statistic/entities"
)

type SalesManagerService interface {
	GetSalesManagerByUserId(userId int32) (*SalesManager, error)
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(sale Sale) (*Sale, error)
	GetSalesManagerGoalByType(period Period, salesManagerId SalesManagerId, typeId SaleTypeId) (SaleAmount, error)
	GetSalesManagerSumsByType(from, to time.Time, salesManagerId SalesManagerId, typeId SaleTypeId) (SaleAmount, error)
	GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error)
	GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]Sale, error)
	GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, period Period) (*[]Sale, error)
	GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error)
	UpdateRatio(smId SalesManagerId, period Period) (Percent, error)
	GetRatio(smId SalesManagerId, period Period) (Percent, error)
	GetSalesManagersOrderedByRatio(pagination Pagination, period Period) (*[]SalesManager, error)
	DeleteSale(saleId SaleId) (*Sale, error)
	EditSale(editBody EditSaleBody) (*Sale, error)
}

type DBSalesManagerService struct {
	repo          repository.SalesManagerRepository
	statisticRepo repository.SalesManagerStatisticRepository
	repository2.SaleTypeRepository
}

func (dbs DBSalesManagerService) EditSale(editBody EditSaleBody) (*Sale, error) {
	return dbs.repo.EditSale(editBody)
}

func (dbs DBSalesManagerService) DeleteSale(saleId SaleId) (*Sale, error) {
	return dbs.repo.DeleteSale(saleId)
}

func (dbs DBSalesManagerService) GetSalesManagersOrderedByRatio(pagination Pagination, period Period) (*[]SalesManager, error) {
	from, to := period.ConvertToTime()
	return dbs.repo.GetSalesManagersListOrderedByRatio(pagination, from, to)
}

func (dbs DBSalesManagerService) GetRatio(smId SalesManagerId, period Period) (Percent, error) {
	from, to := period.ConvertToTime()
	return dbs.statisticRepo.GetSalesManagerRatioByPeriod(smId, from, to)
}

func (dbs DBSalesManagerService) GetManagerSalesByPeriod(salesManagerId SalesManagerId, pagination Pagination, period Period) (*[]Sale, error) {
	from, to := period.ConvertToTime()
	return dbs.repo.GetManagerSalesByPeriod(salesManagerId, pagination, from, to)
}

func (dbs DBSalesManagerService) GetSalesManagerSalesCount(salesManagerId SalesManagerId) (int32, error) {
	return dbs.repo.GetSalesManagerSalesCount(salesManagerId)
}

func (dbs DBSalesManagerService) GetManagerSales(salesManagerId SalesManagerId, pagination Pagination) (*[]Sale, error) {
	return dbs.repo.GetManagerSales(salesManagerId, pagination)
}

func (dbs DBSalesManagerService) SaveSale(sale Sale) (*Sale, error) {
	return dbs.repo.SaveSale(sale.SaleManagerId, sale.SaleDate, sale.SalesAmount, sale.SaleType.Id, sale.SaleDescription)
}

func NewSalesManagerService(repo repository.SalesManagerRepository, statisticRepo repository.SalesManagerStatisticRepository, saleTypeRepo repository2.SaleTypeRepository) SalesManagerService {
	return DBSalesManagerService{
		repo:               repo,
		SaleTypeRepository: saleTypeRepo,
		statisticRepo:      statisticRepo,
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

func (dbs DBSalesManagerService) GetSalesManagerGoalByType(period Period, salesManagerId SalesManagerId, typeId SaleTypeId) (SaleAmount, error) {
	from, to := period.ConvertToTime()
	log.Println(from, to)
	return dbs.statisticRepo.GetSalesGoalBySaleTypeAndManager(salesManagerId, typeId, from, to)
}

func (dbs DBSalesManagerService) GetSalesManagerSumsByType(from, to time.Time, salesManagerId SalesManagerId, typeId SaleTypeId) (SaleAmount, error) {
	return dbs.statisticRepo.GetSalesSumBySaleTypeAndManager(salesManagerId, typeId, from, to)
}

func (dbs DBSalesManagerService) GetSalesManagerYearMonthlyStatistic(smId SalesManagerId, year int32) (*[]MonthlyYearStatistic, error) {
	return dbs.repo.GetMonthlyYearSaleStatistic(smId, year)
}
