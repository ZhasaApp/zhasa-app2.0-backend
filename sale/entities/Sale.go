package entities

import (
	"time"
	"zhasa2.0/manager/entities"
)

type Sale struct {
	Id              SaleId
	SaleManagerId   entities.SalesManagerId
	SalesTypeId     SaleTypeId
	SalesAmount     SaleAmount
	SaleDate        time.Time
	SaleDescription SaleDescription
}

type SaleId int32

type SaleDescription string

type SalesBySaleType map[SaleType]Sale
