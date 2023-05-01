package service

import (
	"zhasa2.0/manager/entities"
	"zhasa2.0/manager/repository"
	sale "zhasa2.0/sale/entities"
	repository2 "zhasa2.0/sale/repository"
)

type SalesManagerService interface {
	GetSalesManagerByUserId(userId int32) (*entities.SalesManager, error)
	CreateSalesManager(userId int32, branchId int32) error
	SaveSale(sale sale.Sale) error
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
	provideManagers() (error, *entities.SalesManagers)
}

func (dbs DBSalesManagerService) ProvideManagers(visitor salesManagerVisitor) (error, *entities.SalesManagers) {
	return visitor.provideManagers()
}

func (dbs DBSalesManagerService) CreateSalesManager(userId int32, branchId int32) error {
	return dbs.repo.CreateSalesManager(userId, branchId)
}

func (dbs DBSalesManagerService) GetSalesManagerByUserId(userId int32) (*entities.SalesManager, error) {
	return dbs.repo.GetSalesManagerByUserId(userId)
}
