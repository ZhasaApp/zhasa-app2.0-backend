package service

import (
	"zhasa2.0/manager/entities"
	"zhasa2.0/manager/repository"
	sale "zhasa2.0/sale/entities"
	"zhasa2.0/statistic"
	entities2 "zhasa2.0/statistic/entities"
)

type SalesManagerService interface {
	SaveSale(sale sale.Sale) error
	statistic.Statistic
}

type DBSalesManagerService struct {
	salesManager entities.SalesManager
	repo         repository.SalesManagerRepository
}

func (dbs DBSalesManagerService) SaveSale(sale sale.Sale) error {
	return dbs.repo.SaveSale(sale.SaleDate, sale.SalesAmount, sale.SalesType.Id)
}

func (dbs DBSalesManagerService) ProvideSalesSums(period entities2.Period) (*statistic.SaleSumByType, error) {
	from, to := period.ConvertToTime()
	return dbs.repo.ProvideSums(from, to)
}
