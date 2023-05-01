package service

import (
	"zhasa2.0/sale/entities"
	"zhasa2.0/sale/repository"
)

type SaleTypeService interface {
	GetSaleType(id entities.SaleTypeId) (*entities.SaleType, error)
}

type DBSaleTypeService struct {
	repo repository.SaleTypeRepository
}

func NewSaleTypeService(repo repository.SaleTypeRepository) SaleTypeService {
	return DBSaleTypeService{
		repo: repo,
	}
}

func (ds DBSaleTypeService) GetSaleType(id entities.SaleTypeId) (*entities.SaleType, error) {
	return ds.repo.GetSaleType(id)
}
