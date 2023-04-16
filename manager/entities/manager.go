package entities

import (
	"time"
	"zhasa2.0/manager/repository"
	sale "zhasa2.0/sale/entities"
)

type SalesManageId int

type SalesManager struct {
	Id   SalesManageId
	repo repository.SalesManagerRepository
}

func (s SalesManager) AddSale(time time.Time, saleType sale.SaleType, amount sale.SaleAmount) {

}
