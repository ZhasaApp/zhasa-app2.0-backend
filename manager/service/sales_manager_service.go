package service

import (
	entities3 "zhasa2.0/branch/entities"
	"zhasa2.0/manager/entities"
	"zhasa2.0/manager/repository"
	sale "zhasa2.0/sale/entities"
	"zhasa2.0/statistic"
	entities2 "zhasa2.0/statistic/entities"
)

type SalesManagerService interface {
	SaveSale(sale sale.Sale) error
	statistic.Statistic
	ProvideManagers(visitor salesManagerVisitor) (error, *entities.SalesManagers)
}
type DBSalesManagerService struct {
	salesManager entities.SalesManager
	repo         repository.SalesManagerRepository
}

type salesManagerVisitor interface {
	provideManagers() (error, *entities.SalesManagers)
}

func (dbs DBSalesManagerService) ProvideManagers(visitor salesManagerVisitor) (error, *entities.SalesManagers) {
	return visitor.provideManagers()
}

type rankedSalesManagersVisitor struct {
	entities.SalesManagers
	entities2.Period
	repo repository.SalesManagerRepository
	page int32
	size int32
}

type branchSalesManagersVisitor struct {
	branchId entities3.BranchId
	repo     repository.SalesManagerRepository
}

func (rv rankedSalesManagersVisitor) provideManagers() (*entities.SalesManagers, error) {
	from, to := rv.Period.ConvertToTime()
	return rv.repo.ProvideRankedSalesManagersList(from, to, rv.size, rv.page)
}

func (bv branchSalesManagersVisitor) provideManagers() (*entities.SalesManagers, error) {
	return bv.repo.ProvideBranchSalesManagers(bv.branchId)
}

func (dbs DBSalesManagerService) SaveSale(sale sale.Sale) error {
	return dbs.repo.SaveSale(sale.SaleDate, sale.SalesAmount, sale.SalesType.Id)
}

func (dbs DBSalesManagerService) ProvideSalesSums(period entities2.Period) (*statistic.SaleSumByType, error) {
	from, to := period.ConvertToTime()
	return dbs.repo.ProvideSums(from, to)
}
